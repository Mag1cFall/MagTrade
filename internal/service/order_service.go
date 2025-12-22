package service

import (
	"context"
	"time"

	"github.com/Mag1cFall/magtrade/internal/cache"
	"github.com/Mag1cFall/magtrade/internal/model"
	"github.com/Mag1cFall/magtrade/internal/mq"
	"github.com/Mag1cFall/magtrade/internal/pkg/utils"
	"github.com/Mag1cFall/magtrade/internal/repository"
	"go.uber.org/zap"
)

type OrderService struct {
	orderRepo     *repository.OrderRepository
	flashSaleRepo *repository.FlashSaleRepository
	stockService  *cache.StockService
	producer      *mq.Producer
	log           *zap.Logger
}

func NewOrderService(producer *mq.Producer, log *zap.Logger) *OrderService {
	return &OrderService{
		orderRepo:     repository.NewOrderRepository(),
		flashSaleRepo: repository.NewFlashSaleRepository(),
		stockService:  cache.NewStockService(),
		producer:      producer,
		log:           log,
	}
}

type OrderListResponse struct {
	Orders   []model.Order `json:"orders"`
	Total    int64         `json:"total"`
	Page     int           `json:"page"`
	PageSize int           `json:"page_size"`
}

func (s *OrderService) CreateFromMessage(ctx context.Context, msg *mq.FlashSaleOrderMessage) (*model.Order, error) {
	flashSale, err := s.flashSaleRepo.GetByID(ctx, msg.FlashSaleID)
	if err != nil {
		return nil, err
	}

	existing, err := s.orderRepo.GetByUserAndFlashSale(ctx, msg.UserID, msg.FlashSaleID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		s.log.Warn("duplicate order creation attempt",
			zap.Int64("user_id", msg.UserID),
			zap.Int64("flash_sale_id", msg.FlashSaleID),
		)
		return existing, nil
	}

	order := &model.Order{
		OrderNo:     utils.GenerateOrderNo(),
		UserID:      msg.UserID,
		FlashSaleID: msg.FlashSaleID,
		Amount:      flashSale.FlashPrice * float64(msg.Quantity),
		Quantity:    msg.Quantity,
		Status:      model.OrderStatusPending,
	}

	if err := s.orderRepo.Create(ctx, order); err != nil {
		return nil, err
	}

	if err := s.flashSaleRepo.DecrementStock(ctx, msg.FlashSaleID, msg.Quantity); err != nil {
		s.log.Error("failed to decrement db stock", zap.Error(err))
	}

	s.log.Info("order created",
		zap.String("order_no", order.OrderNo),
		zap.Int64("user_id", msg.UserID),
		zap.Int64("flash_sale_id", msg.FlashSaleID),
	)

	return order, nil
}

func (s *OrderService) GetByOrderNo(ctx context.Context, userID int64, orderNo string) (*model.Order, error) {
	order, err := s.orderRepo.GetByOrderNo(ctx, orderNo)
	if err != nil {
		return nil, err
	}

	if order.UserID != userID {
		return nil, repository.ErrOrderNotFound
	}

	return order, nil
}

func (s *OrderService) ListByUser(ctx context.Context, userID int64, page, pageSize int) (*OrderListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	orders, total, err := s.orderRepo.ListByUser(ctx, userID, page, pageSize)
	if err != nil {
		return nil, err
	}

	return &OrderListResponse{
		Orders:   orders,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func (s *OrderService) Pay(ctx context.Context, userID int64, orderNo string) (*model.Order, error) {
	order, err := s.orderRepo.GetByOrderNo(ctx, orderNo)
	if err != nil {
		return nil, err
	}

	if order.UserID != userID {
		return nil, repository.ErrOrderNotFound
	}

	if !order.CanPay() {
		return nil, ErrOrderStatusInvalid
	}

	if err := s.orderRepo.Pay(ctx, order.ID); err != nil {
		return nil, err
	}

	s.notifyOrderStatusChange(ctx, order, model.OrderStatusPending, model.OrderStatusPaid)

	order.Status = model.OrderStatusPaid
	now := time.Now()
	order.PaidAt = &now

	return order, nil
}

func (s *OrderService) Cancel(ctx context.Context, userID int64, orderNo string) (*model.Order, error) {
	order, err := s.orderRepo.GetByOrderNo(ctx, orderNo)
	if err != nil {
		return nil, err
	}

	if order.UserID != userID {
		return nil, repository.ErrOrderNotFound
	}

	if !order.CanCancel() {
		return nil, ErrOrderStatusInvalid
	}

	if err := s.orderRepo.Cancel(ctx, order.ID); err != nil {
		return nil, err
	}

	if err := s.stockService.RestoreStock(ctx, order.FlashSaleID, userID, order.Quantity); err != nil {
		s.log.Error("failed to restore redis stock", zap.Error(err))
	}

	if err := s.flashSaleRepo.IncrementStock(ctx, order.FlashSaleID, order.Quantity); err != nil {
		s.log.Error("failed to restore db stock", zap.Error(err))
	}

	s.notifyOrderStatusChange(ctx, order, model.OrderStatusPending, model.OrderStatusCancelled)

	order.Status = model.OrderStatusCancelled
	return order, nil
}

func (s *OrderService) CancelExpiredOrders(ctx context.Context, expireDuration time.Duration) error {
	orders, err := s.orderRepo.CancelExpiredPending(ctx, expireDuration, 100)
	if err != nil {
		return err
	}

	for _, order := range orders {
		if err := s.stockService.RestoreStock(ctx, order.FlashSaleID, order.UserID, order.Quantity); err != nil {
			s.log.Error("failed to restore redis stock for expired order",
				zap.String("order_no", order.OrderNo),
				zap.Error(err),
			)
		}

		if err := s.flashSaleRepo.IncrementStock(ctx, order.FlashSaleID, order.Quantity); err != nil {
			s.log.Error("failed to restore db stock for expired order",
				zap.String("order_no", order.OrderNo),
				zap.Error(err),
			)
		}

		s.notifyOrderStatusChange(ctx, &order, model.OrderStatusPending, model.OrderStatusCancelled)
	}

	if len(orders) > 0 {
		s.log.Info("cancelled expired orders", zap.Int("count", len(orders)))
	}

	return nil
}

func (s *OrderService) notifyOrderStatusChange(ctx context.Context, order *model.Order, oldStatus, newStatus model.OrderStatus) {
	msg := &mq.OrderStatusChangeMessage{
		MessageID: utils.GenerateTicket(),
		Timestamp: time.Now(),
		OrderNo:   order.OrderNo,
		UserID:    order.UserID,
		OldStatus: int(oldStatus),
		NewStatus: int(newStatus),
	}

	if err := s.producer.SendOrderStatusChange(ctx, msg); err != nil {
		s.log.Error("failed to send order status change message", zap.Error(err))
	}
}

var ErrOrderStatusInvalid = repository.ErrOrderNotFound
