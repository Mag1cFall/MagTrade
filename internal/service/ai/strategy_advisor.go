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

const strategyAnalysisPrompt = `你是一个秒杀活动分析专家。根据提供的活动数据，分析活动热度并给出策略建议。

请严格按照以下JSON格式返回分析结果：
{
  "difficulty_score": <1-10的整数，10表示最难抢>,
  "difficulty_reason": "<难度评估原因>",
  "timing_advice": "<最佳抢购时机建议>",
  "success_probability": <0-1的小数，预估成功概率>,
  "recommendations": ["<建议1>", "<建议2>", "<建议3>"]
}

只返回JSON，不要有其他内容。`

type StrategyAdvisor struct {
	llm           *LLMClient
	flashSaleRepo *repository.FlashSaleRepository
	recRepo       *repository.AIRecommendationRepository
	stockService  *cache.StockService
	log           *zap.Logger
}

func NewStrategyAdvisor(llm *LLMClient, log *zap.Logger) *StrategyAdvisor {
	return &StrategyAdvisor{
		llm:           llm,
		flashSaleRepo: repository.NewFlashSaleRepository(),
		recRepo:       repository.NewAIRecommendationRepository(),
		stockService:  cache.NewStockService(),
		log:           log,
	}
}

type StrategyAnalysis struct {
	DifficultyScore    int      `json:"difficulty_score"`
	DifficultyReason   string   `json:"difficulty_reason"`
	TimingAdvice       string   `json:"timing_advice"`
	SuccessProbability float64  `json:"success_probability"`
	Recommendations    []string `json:"recommendations"`
}

type StrategyRecommendation struct {
	FlashSaleID int64             `json:"flash_sale_id"`
	Analysis    *StrategyAnalysis `json:"analysis"`
}

func (s *StrategyAdvisor) AnalyzeFlashSale(ctx context.Context, flashSaleID int64) (*StrategyRecommendation, error) {
	flashSale, err := s.flashSaleRepo.GetByID(ctx, flashSaleID)
	if err != nil {
		return nil, err
	}

	currentStock, err := s.stockService.GetStock(ctx, flashSaleID)
	if err != nil {
		currentStock = flashSale.AvailableStock
	}

	activityData := map[string]interface{}{
		"product_name":   flashSale.Product.Name,
		"original_price": flashSale.Product.OriginalPrice,
		"flash_price":    flashSale.FlashPrice,
		"discount_rate":  (flashSale.Product.OriginalPrice - flashSale.FlashPrice) / flashSale.Product.OriginalPrice * 100,
		"total_stock":    flashSale.TotalStock,
		"current_stock":  currentStock,
		"sold_count":     flashSale.TotalStock - currentStock,
		"per_user_limit": flashSale.PerUserLimit,
		"start_time":     flashSale.StartTime.Format("2006-01-02 15:04:05"),
		"end_time":       flashSale.EndTime.Format("2006-01-02 15:04:05"),
		"status":         flashSale.Status,
	}

	dataJSON, _ := json.Marshal(activityData)
	userMessage := fmt.Sprintf("请分析以下秒杀活动：\n%s", string(dataJSON))

	response, err := s.llm.ChatWithSystem(ctx, strategyAnalysisPrompt, userMessage)
	if err != nil {
		return nil, fmt.Errorf("failed to get LLM analysis: %w", err)
	}

	var analysis StrategyAnalysis
	if err := json.Unmarshal([]byte(response), &analysis); err != nil {
		s.log.Warn("failed to parse LLM response as JSON, using fallback",
			zap.String("response", response),
			zap.Error(err),
		)

		analysis = s.fallbackAnalysis(flashSale, currentStock)
	}

	rec := &model.AIRecommendation{
		FlashSaleID:        flashSaleID,
		RecommendationType: model.RecommendationTypeTimingAdvice,
		Content:            response,
		ConfidenceScore:    0.8,
	}
	if err := s.recRepo.Create(ctx, rec); err != nil {
		s.log.Error("failed to save recommendation", zap.Error(err))
	}

	return &StrategyRecommendation{
		FlashSaleID: flashSaleID,
		Analysis:    &analysis,
	}, nil
}

func (s *StrategyAdvisor) fallbackAnalysis(flashSale *model.FlashSale, currentStock int) StrategyAnalysis {
	discountRate := (flashSale.Product.OriginalPrice - flashSale.FlashPrice) / flashSale.Product.OriginalPrice * 100

	difficultyScore := 5
	if discountRate > 70 {
		difficultyScore += 2
	} else if discountRate > 50 {
		difficultyScore += 1
	}

	if flashSale.TotalStock < 100 {
		difficultyScore += 2
	} else if flashSale.TotalStock < 500 {
		difficultyScore += 1
	}

	if difficultyScore > 10 {
		difficultyScore = 10
	}

	successProb := float64(currentStock) / float64(flashSale.TotalStock) * 0.5
	if successProb > 0.5 {
		successProb = 0.5
	}

	return StrategyAnalysis{
		DifficultyScore:    difficultyScore,
		DifficultyReason:   fmt.Sprintf("折扣力度%.0f%%，库存%d件", discountRate, flashSale.TotalStock),
		TimingAdvice:       "建议在活动开始前30秒进入页面，开始后立即点击",
		SuccessProbability: successProb,
		Recommendations: []string{
			"确保网络稳定，建议使用有线网络",
			"提前登录并填写好收货地址",
			"建议使用多设备同时抢购",
		},
	}
}

func (s *StrategyAdvisor) GetLatestRecommendation(ctx context.Context, flashSaleID int64) (*model.AIRecommendation, error) {
	return s.recRepo.GetLatestByType(ctx, flashSaleID, model.RecommendationTypeTimingAdvice)
}
