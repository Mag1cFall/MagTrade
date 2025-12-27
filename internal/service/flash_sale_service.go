// 秒殺活動業務邏輯服務
//
// 本檔案是整個秒殺系統的核心，包含：
// - 活動建立、查詢、狀態管理
// - Rush 方法：秒殺搶購核心流程
// 流程：驗證 → 限購檢查 → 分散式鎖 → Redis 扣減 → Kafka 非同步下單
package service

import (
	"context"
	"errors"
	"time"

	"github.com/Mag1cFall/magtrade/internal/cache"
	"github.com/Mag1cFall/magtrade/internal/model"
	"github.com/Mag1cFall/magtrade/internal/mq"
	"github.com/Mag1cFall/magtrade/internal/pkg/utils"
	"github.com/Mag1cFall/magtrade/internal/repository"
	"go.uber.org/zap"
)

// 業務錯誤定義
var (
	ErrFlashSaleNotActive  = errors.New("flash sale is not active")
	ErrStockInsufficient   = errors.New("stock insufficient")
	ErrLimitExceeded       = errors.New("purchase limit exceeded")
	ErrAlreadyPurchased    = errors.New("already purchased in this flash sale")
	ErrFlashSaleNotStarted = errors.New("flash sale has not started")
	ErrFlashSaleEnded      = errors.New("flash sale has ended")
)

// FlashSaleService 秒殺業務服務
type FlashSaleService struct {
	flashSaleRepo *repository.FlashSaleRepository
	orderRepo     *repository.OrderRepository
	stockService  *cache.StockService // Redis 庫存操作
	producer      *mq.Producer        // Kafka 訊息生產者
	log           *zap.Logger
}

func NewFlashSaleService(producer *mq.Producer, log *zap.Logger) *FlashSaleService {
	return &FlashSaleService{
		flashSaleRepo: repository.NewFlashSaleRepository(),
		orderRepo:     repository.NewOrderRepository(),
		stockService:  cache.NewStockService(),
		producer:      producer,
		log:           log,
	}
}

// CreateFlashSaleRequest 建立秒殺活動請求
type CreateFlashSaleRequest struct {
	ProductID    int64   `json:"product_id" binding:"required"`
	FlashPrice   float64 `json:"flash_price" binding:"required,gt=0"`
	TotalStock   int     `json:"total_stock" binding:"required,gt=0"`
	PerUserLimit int     `json:"per_user_limit" binding:"omitempty,gt=0"`
	StartTime    string  `json:"start_time" binding:"required"` // RFC3339 格式
	EndTime      string  `json:"end_time" binding:"required"`
}

// FlashSaleDetailResponse 活動詳情回應（含即時庫存）
type FlashSaleDetailResponse struct {
	FlashSale    *model.FlashSale `json:"flash_sale"`
	CurrentStock int              `json:"current_stock"` // 從 Redis 取得的即時庫存
	ServerTime   time.Time        `json:"server_time"`   // 伺服器時間，用於前端倒數同步
}

// RushRequest 秒殺搶購請求
type RushRequest struct {
	Quantity int `json:"quantity" binding:"omitempty,gt=0"`
}

// RushResponse 秒殺搶購回應
type RushResponse struct {
	Success  bool   `json:"success"`
	Ticket   string `json:"ticket,omitempty"`   // 排隊憑證，用於查詢訂單狀態
	Position int    `json:"position,omitempty"` // 排隊位置（保留欄位）
	Message  string `json:"message"`
	OrderNo  string `json:"order_no,omitempty"` // 若已購買過則返回訂單號
}

// Create 建立秒殺活動
func (s *FlashSaleService) Create(ctx context.Context, req *CreateFlashSaleRequest) (*model.FlashSale, error) {
	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		return nil, errors.New("invalid start_time format, use RFC3339")
	}

	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		return nil, errors.New("invalid end_time format, use RFC3339")
	}

	if endTime.Before(startTime) {
		return nil, errors.New("end_time must be after start_time")
	}

	perUserLimit := req.PerUserLimit
	if perUserLimit <= 0 {
		perUserLimit = 1 // 預設每人限購 1 件
	}

	flashSale := &model.FlashSale{
		ProductID:      req.ProductID,
		FlashPrice:     req.FlashPrice,
		TotalStock:     req.TotalStock,
		AvailableStock: req.TotalStock,
		PerUserLimit:   perUserLimit,
		StartTime:      startTime,
		EndTime:        endTime,
		Status:         model.FlashSaleStatusPending,
	}

	if err := s.flashSaleRepo.Create(ctx, flashSale); err != nil {
		return nil, err
	}

	// 同步初始化 Redis 庫存
	if err := s.stockService.InitStock(ctx, flashSale.ID, req.TotalStock); err != nil {
		s.log.Error("failed to init redis stock", zap.Int64("flash_sale_id", flashSale.ID), zap.Error(err))
	}

	return flashSale, nil
}

// GetByID 查詢活動詳情（含即時庫存）
func (s *FlashSaleService) GetByID(ctx context.Context, id int64) (*FlashSaleDetailResponse, error) {
	flashSale, err := s.flashSaleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	stock, err := s.stockService.GetStock(ctx, id)
	if err != nil {
		stock = flashSale.AvailableStock // Redis 失敗時降級使用 DB 庫存
	}

	return &FlashSaleDetailResponse{
		FlashSale:    flashSale,
		CurrentStock: stock,
		ServerTime:   time.Now(),
	}, nil
}

// List 分頁查詢活動列表
func (s *FlashSaleService) List(ctx context.Context, page, pageSize int, status *model.FlashSaleStatus) ([]model.FlashSale, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	return s.flashSaleRepo.List(ctx, page, pageSize, status)
}

// ListActive 查詢進行中的活動
func (s *FlashSaleService) ListActive(ctx context.Context) ([]model.FlashSale, error) {
	return s.flashSaleRepo.ListActive(ctx)
}

// GetStock 查詢即時庫存
func (s *FlashSaleService) GetStock(ctx context.Context, id int64) (int, error) {
	return s.stockService.GetStock(ctx, id)
}

// Rush 秒殺搶購核心邏輯
// 這是整個系統最關鍵的方法，包含完整的併發控制流程
func (s *FlashSaleService) Rush(ctx context.Context, userID, flashSaleID int64, quantity int) (*RushResponse, error) {
	if quantity <= 0 {
		quantity = 1
	}

	// 階段一：驗證活動狀態
	flashSale, err := s.flashSaleRepo.GetByID(ctx, flashSaleID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	if now.Before(flashSale.StartTime) {
		return &RushResponse{
			Success: false,
			Message: "秒杀活动尚未开始",
		}, ErrFlashSaleNotStarted
	}
	if now.After(flashSale.EndTime) {
		return &RushResponse{
			Success: false,
			Message: "秒杀活动已结束",
		}, ErrFlashSaleEnded
	}
	if flashSale.Status != model.FlashSaleStatusActive {
		return &RushResponse{
			Success: false,
			Message: "秒杀活动未开放",
		}, ErrFlashSaleNotActive
	}

	// 階段二：檢查使用者是否已購買
	existingOrder, err := s.orderRepo.GetByUserAndFlashSale(ctx, userID, flashSaleID)
	if err != nil {
		return nil, err
	}
	if existingOrder != nil {
		return &RushResponse{
			Success: false,
			Message: "您已参与过本次秒杀",
			OrderNo: existingOrder.OrderNo,
		}, ErrAlreadyPurchased
	}

	// 階段三：取得分散式鎖（防止同一使用者併發重複提交）
	lock := cache.NewDistributedLock(flashSaleID, userID)
	acquired, err := lock.Lock(ctx, 10000) // 鎖 10 秒
	if err != nil {
		s.log.Error("failed to acquire lock", zap.Error(err))
		return nil, errors.New("system busy, please retry")
	}
	if !acquired {
		return &RushResponse{
			Success: false,
			Message: "请勿重复提交",
		}, nil
	}
	defer func() {
		if err := lock.Unlock(ctx); err != nil {
			s.log.Error("failed to release lock", zap.Error(err))
		}
	}()

	// 階段四：Redis 原子扣減庫存
	result, err := s.stockService.DeductStock(ctx, flashSaleID, userID, quantity, flashSale.PerUserLimit)
	if err != nil {
		return nil, err
	}

	if !result.Success {
		switch result.Code {
		case -1: // 庫存不足
			return &RushResponse{Success: false, Message: result.Message}, ErrStockInsufficient
		case -2: // 超過限購
			return &RushResponse{Success: false, Message: result.Message}, ErrLimitExceeded
		default:
			return &RushResponse{Success: false, Message: result.Message}, nil
		}
	}

	// 階段五：發送 Kafka 非同步訂單消息
	ticket := utils.GenerateTicket() // 生成排隊憑證

	msg := &mq.FlashSaleOrderMessage{
		MessageID:   ticket,
		Timestamp:   time.Now(),
		FlashSaleID: flashSaleID,
		UserID:      userID,
		Quantity:    quantity,
		Ticket:      ticket,
	}

	if err := s.producer.SendFlashSaleOrder(ctx, msg); err != nil {
		// 發送失敗，必須回滾 Redis 庫存
		s.log.Error("failed to send order message, restoring stock",
			zap.Error(err),
			zap.Int64("user_id", userID),
			zap.Int64("flash_sale_id", flashSaleID),
		)

		if restoreErr := s.stockService.RestoreStock(ctx, flashSaleID, userID, quantity); restoreErr != nil {
			s.log.Error("failed to restore stock", zap.Error(restoreErr))
		}

		return nil, errors.New("system busy, please retry")
	}

	// 搶購成功，返回排隊憑證
	return &RushResponse{
		Success:  true,
		Ticket:   ticket,
		Position: 0,
		Message:  "排队中，请等待结果",
	}, nil
}

// ActivatePendingFlashSales 自動開啟已到時間的待開始活動
func (s *FlashSaleService) ActivatePendingFlashSales(ctx context.Context) error {
	affected, err := s.flashSaleRepo.UpdatePendingToActive(ctx)
	if err != nil {
		return err
	}
	if affected > 0 {
		s.log.Info("activated flash sales", zap.Int64("count", affected))
	}
	return nil
}

// FinishExpiredFlashSales 自動結束已過期的進行中活動
func (s *FlashSaleService) FinishExpiredFlashSales(ctx context.Context) error {
	affected, err := s.flashSaleRepo.UpdateActiveToFinished(ctx)
	if err != nil {
		return err
	}
	if affected > 0 {
		s.log.Info("finished flash sales", zap.Int64("count", affected))
	}
	return nil
}
