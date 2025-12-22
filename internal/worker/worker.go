package worker

import (
	"context"
	"time"

	"github.com/Mag1cFall/magtrade/internal/handler"
	"github.com/Mag1cFall/magtrade/internal/mq"
	"github.com/Mag1cFall/magtrade/internal/service"
	"go.uber.org/zap"
)

type OrderWorker struct {
	orderService *service.OrderService
	wsHub        *handler.WSHub
	log          *zap.Logger
}

func NewOrderWorker(producer *mq.Producer, wsHub *handler.WSHub, log *zap.Logger) *OrderWorker {
	return &OrderWorker{
		orderService: service.NewOrderService(producer, log),
		wsHub:        wsHub,
		log:          log,
	}
}

func (w *OrderWorker) HandleFlashSaleOrder(ctx context.Context, data []byte) error {
	msg, err := mq.ParseFlashSaleOrderMessage(data)
	if err != nil {
		w.log.Error("failed to parse flash sale order message", zap.Error(err))
		return err
	}

	order, err := w.orderService.CreateFromMessage(ctx, msg)
	if err != nil {
		w.log.Error("failed to create order",
			zap.Error(err),
			zap.Int64("user_id", msg.UserID),
			zap.Int64("flash_sale_id", msg.FlashSaleID),
		)

		w.wsHub.SendToUser(msg.UserID, "flash_sale_result", map[string]interface{}{
			"flash_sale_id": msg.FlashSaleID,
			"success":       false,
			"message":       "订单创建失败，请稍后查看",
			"ticket":        msg.Ticket,
		})

		return err
	}

	w.wsHub.SendToUser(msg.UserID, "flash_sale_result", map[string]interface{}{
		"flash_sale_id": msg.FlashSaleID,
		"success":       true,
		"order_no":      order.OrderNo,
		"message":       "恭喜您抢购成功！",
		"ticket":        msg.Ticket,
	})

	w.log.Info("order created and notified",
		zap.String("order_no", order.OrderNo),
		zap.Int64("user_id", msg.UserID),
	)

	return nil
}

func (w *OrderWorker) HandleOrderStatusChange(ctx context.Context, data []byte) error {
	msg, err := mq.ParseOrderStatusChangeMessage(data)
	if err != nil {
		w.log.Error("failed to parse order status change message", zap.Error(err))
		return err
	}

	w.wsHub.SendToUser(msg.UserID, "order_status_change", map[string]interface{}{
		"order_no":   msg.OrderNo,
		"old_status": msg.OldStatus,
		"new_status": msg.NewStatus,
	})

	return nil
}

type SchedulerWorker struct {
	flashSaleService *service.FlashSaleService
	orderService     *service.OrderService
	log              *zap.Logger
	stopCh           chan struct{}
}

func NewSchedulerWorker(producer *mq.Producer, log *zap.Logger) *SchedulerWorker {
	return &SchedulerWorker{
		flashSaleService: service.NewFlashSaleService(producer, log),
		orderService:     service.NewOrderService(producer, log),
		log:              log,
		stopCh:           make(chan struct{}),
	}
}

func (w *SchedulerWorker) Start(ctx context.Context) {
	go w.runFlashSaleStatusUpdater(ctx)
	go w.runExpiredOrderCanceller(ctx)
}

func (w *SchedulerWorker) Stop() {
	close(w.stopCh)
}

func (w *SchedulerWorker) runFlashSaleStatusUpdater(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-w.stopCh:
			return
		case <-ticker.C:
			if err := w.flashSaleService.ActivatePendingFlashSales(ctx); err != nil {
				w.log.Error("failed to activate pending flash sales", zap.Error(err))
			}

			if err := w.flashSaleService.FinishExpiredFlashSales(ctx); err != nil {
				w.log.Error("failed to finish expired flash sales", zap.Error(err))
			}
		}
	}
}

func (w *SchedulerWorker) runExpiredOrderCanceller(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	expireDuration := 15 * time.Minute

	for {
		select {
		case <-ctx.Done():
			return
		case <-w.stopCh:
			return
		case <-ticker.C:
			if err := w.orderService.CancelExpiredOrders(ctx, expireDuration); err != nil {
				w.log.Error("failed to cancel expired orders", zap.Error(err))
			}
		}
	}
}
