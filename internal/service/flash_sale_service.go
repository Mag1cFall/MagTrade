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

var (
	ErrFlashSaleNotActive  = errors.New("flash sale is not active")
	ErrStockInsufficient   = errors.New("stock insufficient")
	ErrLimitExceeded       = errors.New("purchase limit exceeded")
	ErrAlreadyPurchased    = errors.New("already purchased in this flash sale")
	ErrFlashSaleNotStarted = errors.New("flash sale has not started")
	ErrFlashSaleEnded      = errors.New("flash sale has ended")
)

type FlashSaleService struct {
	flashSaleRepo *repository.FlashSaleRepository
	orderRepo     *repository.OrderRepository
	stockService  *cache.StockService
	producer      *mq.Producer
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

type CreateFlashSaleRequest struct {
	ProductID    int64   `json:"product_id" binding:"required"`
	FlashPrice   float64 `json:"flash_price" binding:"required,gt=0"`
	TotalStock   int     `json:"total_stock" binding:"required,gt=0"`
	PerUserLimit int     `json:"per_user_limit" binding:"omitempty,gt=0"`
	StartTime    string  `json:"start_time" binding:"required"`
	EndTime      string  `json:"end_time" binding:"required"`
}

type FlashSaleDetailResponse struct {
	FlashSale    *model.FlashSale `json:"flash_sale"`
	CurrentStock int              `json:"current_stock"`
	ServerTime   time.Time        `json:"server_time"`
}

type RushRequest struct {
	Quantity int `json:"quantity" binding:"omitempty,gt=0"`
}

type RushResponse struct {
	Success  bool   `json:"success"`
	Ticket   string `json:"ticket,omitempty"`
	Position int    `json:"position,omitempty"`
	Message  string `json:"message"`
	OrderNo  string `json:"order_no,omitempty"`
}

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
		perUserLimit = 1
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

	if err := s.stockService.InitStock(ctx, flashSale.ID, req.TotalStock); err != nil {
		s.log.Error("failed to init redis stock", zap.Int64("flash_sale_id", flashSale.ID), zap.Error(err))
	}

	return flashSale, nil
}

func (s *FlashSaleService) GetByID(ctx context.Context, id int64) (*FlashSaleDetailResponse, error) {
	flashSale, err := s.flashSaleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	stock, err := s.stockService.GetStock(ctx, id)
	if err != nil {
		stock = flashSale.AvailableStock
	}

	return &FlashSaleDetailResponse{
		FlashSale:    flashSale,
		CurrentStock: stock,
		ServerTime:   time.Now(),
	}, nil
}

func (s *FlashSaleService) List(ctx context.Context, page, pageSize int, status *model.FlashSaleStatus) ([]model.FlashSale, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	return s.flashSaleRepo.List(ctx, page, pageSize, status)
}

func (s *FlashSaleService) ListActive(ctx context.Context) ([]model.FlashSale, error) {
	return s.flashSaleRepo.ListActive(ctx)
}

func (s *FlashSaleService) GetStock(ctx context.Context, id int64) (int, error) {
	return s.stockService.GetStock(ctx, id)
}

func (s *FlashSaleService) Rush(ctx context.Context, userID, flashSaleID int64, quantity int) (*RushResponse, error) {
	if quantity <= 0 {
		quantity = 1
	}

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

	lock := cache.NewDistributedLock(flashSaleID, userID)
	acquired, err := lock.Lock(ctx, 10000)
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

	result, err := s.stockService.DeductStock(ctx, flashSaleID, userID, quantity, flashSale.PerUserLimit)
	if err != nil {
		return nil, err
	}

	if !result.Success {
		switch result.Code {
		case -1:
			return &RushResponse{Success: false, Message: result.Message}, ErrStockInsufficient
		case -2:
			return &RushResponse{Success: false, Message: result.Message}, ErrLimitExceeded
		default:
			return &RushResponse{Success: false, Message: result.Message}, nil
		}
	}

	ticket := utils.GenerateTicket()

	msg := &mq.FlashSaleOrderMessage{
		MessageID:   ticket,
		Timestamp:   time.Now(),
		FlashSaleID: flashSaleID,
		UserID:      userID,
		Quantity:    quantity,
		Ticket:      ticket,
	}

	if err := s.producer.SendFlashSaleOrder(ctx, msg); err != nil {
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

	return &RushResponse{
		Success:  true,
		Ticket:   ticket,
		Position: 0,
		Message:  "排队中，请等待结果",
	}, nil
}

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
