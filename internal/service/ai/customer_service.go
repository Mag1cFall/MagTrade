// 智慧客服服務
//
// 本檔案提供 AI 智慧客服功能
// 整合 LLM、對話歷史、秒殺活動和訂單資料
// 支援上下文感知對話（自動帶入相關業務資料）
package ai

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Mag1cFall/magtrade/internal/cache"
	"github.com/Mag1cFall/magtrade/internal/model"
	"github.com/Mag1cFall/magtrade/internal/repository"
	"go.uber.org/zap"
)

// 客服系統提示詞
const customerServiceSystemPrompt = `你是 MagTrade 秒杀平台的智能客服助手。你的职责是：
1. 回答用户关于秒杀活动、商品信息的问题
2. 查询用户的订单状态
3. 提供秒杀技巧和建议

当用户询问具体活动信息时，你会收到活动的实时数据。请根据这些数据准确回答。
回答时请简洁明了，使用友好的语气。如涉及时间请使用北京时间。
如果用户的问题超出你的能力范围，请礼貌地告知并建议联系人工客服。`

// CustomerService 智慧客服服務
type CustomerService struct {
	llm           *LLMClient
	chatRepo      *repository.ChatHistoryRepository
	flashSaleRepo *repository.FlashSaleRepository
	orderRepo     *repository.OrderRepository
	stockService  *cache.StockService
	log           *zap.Logger
}

func NewCustomerService(llm *LLMClient, log *zap.Logger) *CustomerService {
	return &CustomerService{
		llm:           llm,
		chatRepo:      repository.NewChatHistoryRepository(),
		flashSaleRepo: repository.NewFlashSaleRepository(),
		orderRepo:     repository.NewOrderRepository(),
		stockService:  cache.NewStockService(),
		log:           log,
	}
}

// ChatRequest 對話請求
type ChatRequest struct {
	SessionID string `json:"session_id" binding:"required"` // 對話 Session ID
	Message   string `json:"message" binding:"required,max=2000"`
}

// ChatResponse 對話回應
type ChatResponse struct {
	SessionID   string                 `json:"session_id"`
	Response    string                 `json:"response"`
	RelatedData map[string]interface{} `json:"related_data,omitempty"` // 相關業務資料
}

// Chat 普通對話（等待完整回應）
func (s *CustomerService) Chat(ctx context.Context, userID int64, req *ChatRequest) (*ChatResponse, error) {
	// 取得歷史對話（最近 10 條）
	history, err := s.chatRepo.GetBySession(ctx, userID, req.SessionID, 10)
	if err != nil {
		s.log.Error("failed to get chat history", zap.Error(err))
	}

	// 蒐集上下文（活動資訊、訂單資訊）
	contextInfo, relatedData := s.gatherContext(ctx, userID, req.Message)

	// 組合系統提示詞
	enhancedSystemPrompt := customerServiceSystemPrompt
	if contextInfo != "" {
		enhancedSystemPrompt += "\n\n当前上下文信息：\n" + contextInfo
	}

	// 轉換歷史格式
	var chatHistory []ChatMessage
	for _, h := range history {
		chatHistory = append(chatHistory, ChatMessage{
			Role:    string(h.Role),
			Content: h.Content,
		})
	}

	// 呼叫 LLM
	response, err := s.llm.ChatWithHistory(ctx, enhancedSystemPrompt, chatHistory, req.Message)
	if err != nil {
		return nil, fmt.Errorf("failed to get LLM response: %w", err)
	}

	// 儲存使用者消息
	userChat := &model.ChatHistory{
		UserID:    userID,
		SessionID: req.SessionID,
		Role:      model.ChatRoleUser,
		Content:   req.Message,
	}
	if err := s.chatRepo.Create(ctx, userChat); err != nil {
		s.log.Error("failed to save user message", zap.Error(err))
	}

	// 儲存 AI 回應
	assistantChat := &model.ChatHistory{
		UserID:    userID,
		SessionID: req.SessionID,
		Role:      model.ChatRoleAssistant,
		Content:   response,
	}
	if err := s.chatRepo.Create(ctx, assistantChat); err != nil {
		s.log.Error("failed to save assistant message", zap.Error(err))
	}

	return &ChatResponse{
		SessionID:   req.SessionID,
		Response:    response,
		RelatedData: relatedData,
	}, nil
}

// ChatStream 串流對話（逐字回傳）
func (s *CustomerService) ChatStream(ctx context.Context, userID int64, req *ChatRequest, chunkHandler func(chunk StreamChunk) error) error {
	history, err := s.chatRepo.GetBySession(ctx, userID, req.SessionID, 10)
	if err != nil {
		s.log.Error("failed to get chat history", zap.Error(err))
	}

	contextInfo, _ := s.gatherContext(ctx, userID, req.Message)

	enhancedSystemPrompt := customerServiceSystemPrompt
	if contextInfo != "" {
		enhancedSystemPrompt += "\n\n当前上下文信息：\n" + contextInfo
	}

	var chatHistory []ChatMessage
	for _, h := range history {
		chatHistory = append(chatHistory, ChatMessage{
			Role:    string(h.Role),
			Content: h.Content,
		})
	}

	// 儲存使用者消息
	userChat := &model.ChatHistory{
		UserID:    userID,
		SessionID: req.SessionID,
		Role:      model.ChatRoleUser,
		Content:   req.Message,
	}
	if err := s.chatRepo.Create(ctx, userChat); err != nil {
		s.log.Error("failed to save user message", zap.Error(err))
	}

	// 串流呼叫並累積回應
	var fullResponse string
	err = s.llm.ChatStreamWithHistory(ctx, enhancedSystemPrompt, chatHistory, req.Message, func(chunk StreamChunk) error {
		if chunk.Type == "content" {
			fullResponse += chunk.Content
		}
		return chunkHandler(chunk)
	})

	if err != nil {
		return fmt.Errorf("failed to get LLM response: %w", err)
	}

	// 儲存完整的 AI 回應
	assistantChat := &model.ChatHistory{
		UserID:    userID,
		SessionID: req.SessionID,
		Role:      model.ChatRoleAssistant,
		Content:   fullResponse,
	}
	if err := s.chatRepo.Create(ctx, assistantChat); err != nil {
		s.log.Error("failed to save assistant message", zap.Error(err))
	}

	return nil
}

// gatherContext 蒐集上下文資訊（活動、訂單）
func (s *CustomerService) gatherContext(ctx context.Context, userID int64, message string) (string, map[string]interface{}) {
	var contextParts []string
	relatedData := make(map[string]interface{})

	// 進行中的秒殺活動
	activeFlashSales, err := s.flashSaleRepo.ListActive(ctx)
	if err == nil && len(activeFlashSales) > 0 {
		var salesInfo []map[string]interface{}
		for _, fs := range activeFlashSales {
			stock, _ := s.stockService.GetStock(ctx, fs.ID)
			info := map[string]interface{}{
				"id":         fs.ID,
				"name":       fs.Product.Name,
				"price":      fs.FlashPrice,
				"stock":      stock,
				"start_time": fs.StartTime.Format("2006-01-02 15:04"),
				"end_time":   fs.EndTime.Format("2006-01-02 15:04"),
			}
			salesInfo = append(salesInfo, info)
		}

		data, _ := json.Marshal(salesInfo)
		contextParts = append(contextParts, fmt.Sprintf("当前进行中的秒杀活动：%s", string(data)))
		relatedData["active_flash_sales"] = salesInfo
	}

	// 即將開始的秒殺活動
	upcomingFlashSales, err := s.flashSaleRepo.ListUpcoming(ctx, 5)
	if err == nil && len(upcomingFlashSales) > 0 {
		var salesInfo []map[string]interface{}
		for _, fs := range upcomingFlashSales {
			info := map[string]interface{}{
				"id":         fs.ID,
				"name":       fs.Product.Name,
				"price":      fs.FlashPrice,
				"stock":      fs.TotalStock,
				"start_time": fs.StartTime.Format("2006-01-02 15:04"),
			}
			salesInfo = append(salesInfo, info)
		}

		data, _ := json.Marshal(salesInfo)
		contextParts = append(contextParts, fmt.Sprintf("即将开始的秒杀活动：%s", string(data)))
		relatedData["upcoming_flash_sales"] = salesInfo
	}

	// 使用者最近訂單
	orders, _, err := s.orderRepo.ListByUser(ctx, userID, 1, 5)
	if err == nil && len(orders) > 0 {
		var ordersInfo []map[string]interface{}
		for _, o := range orders {
			info := map[string]interface{}{
				"order_no": o.OrderNo,
				"amount":   o.Amount,
				"status":   o.Status.String(),
				"time":     o.CreatedAt.Format("2006-01-02 15:04"),
			}
			ordersInfo = append(ordersInfo, info)
		}

		data, _ := json.Marshal(ordersInfo)
		contextParts = append(contextParts, fmt.Sprintf("用户最近的订单：%s", string(data)))
		relatedData["recent_orders"] = ordersInfo
	}

	contextStr := ""
	for _, part := range contextParts {
		contextStr += part + "\n"
	}

	return contextStr, relatedData
}

// GetHistory 取得對話歷史
func (s *CustomerService) GetHistory(ctx context.Context, userID int64, sessionID string) ([]model.ChatHistory, error) {
	return s.chatRepo.GetBySession(ctx, userID, sessionID, 50)
}

// ClearHistory 清除對話歷史
func (s *CustomerService) ClearHistory(ctx context.Context, userID int64, sessionID string) error {
	return s.chatRepo.DeleteBySession(ctx, userID, sessionID)
}
