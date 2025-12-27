// 背景工作者
//
// 本檔案定義消費 Kafka 訊息和定時任務的 Worker
// OrderWorker：處理訂單建立訊息，通過 WebSocket 通知使用者
// SchedulerWorker：定時更新活動狀態、取消過期訂單
package worker

import (
	"context"
	"time"

	"github.com/Mag1cFall/magtrade/internal/handler"
	"github.com/Mag1cFall/magtrade/internal/mq"
	"github.com/Mag1cFall/magtrade/internal/service"
	"go.uber.org/zap"
)

// OrderWorker 訂單處理工作者
// 消費 Kafka 訊息，建立訂單並通過 WebSocket 通知使用者結果
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

// HandleFlashSaleOrder 處理秒殺訂單訊息
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

		// 通知使用者失敗
		w.wsHub.SendToUser(msg.UserID, "flash_sale_result", map[string]interface{}{
			"flash_sale_id": msg.FlashSaleID,
			"success":       false,
			"message":       "订单创建失败，请稍后查看",
			"ticket":        msg.Ticket,
		})

		return err
	}

	// 通知使用者成功
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

// HandleOrderStatusChange 處理訂單狀態變更訊息
func (w *OrderWorker) HandleOrderStatusChange(ctx context.Context, data []byte) error {
	msg, err := mq.ParseOrderStatusChangeMessage(data)
	if err != nil {
		w.log.Error("failed to parse order status change message", zap.Error(err))
		return err
	}

	// 通知使用者狀態變更
	w.wsHub.SendToUser(msg.UserID, "order_status_change", map[string]interface{}{
		"order_no":   msg.OrderNo,
		"old_status": msg.OldStatus,
		"new_status": msg.NewStatus,
	})

	return nil
}

// SchedulerWorker 定時任務工作者
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

// Start 啟動定時任務
func (w *SchedulerWorker) Start(ctx context.Context) {
	go w.runFlashSaleStatusUpdater(ctx)
	go w.runExpiredOrderCanceller(ctx)
}

// Stop 停止定時任務
func (w *SchedulerWorker) Stop() {
	close(w.stopCh)
}

// runFlashSaleStatusUpdater 定時更新秒殺活動狀態
// 待開始 → 進行中、進行中 → 已結束
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

// runExpiredOrderCanceller 定時取消過期未付款訂單
// 15 分鐘未付款的訂單自動取消並恢復庫存
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
