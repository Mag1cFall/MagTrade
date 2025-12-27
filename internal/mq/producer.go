// Kafka 訊息生產者
//
// 本檔案封裝 Kafka 生產者功能，用於發送非同步訊息
// 主要訊息類型：秒殺訂單、訂單狀態變更、AI 分析任務
// 使用 segmentio/kafka-go 客戶端，支援批次發送
package mq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Mag1cFall/magtrade/internal/config"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

// Producer Kafka 生產者
type Producer struct {
	writers map[string]*kafka.Writer // 每個 Topic 一個 Writer
	log     *zap.Logger
}

// NewProducer 建立生產者，為每個 Topic 初始化獨立的 Writer
func NewProducer(cfg *config.KafkaConfig, log *zap.Logger) *Producer {
	writers := make(map[string]*kafka.Writer)

	topics := []string{
		cfg.Topics.FlashSaleOrders,
		cfg.Topics.OrderStatusChange,
		cfg.Topics.AIAnalysisTasks,
	}

	for _, topic := range topics {
		writers[topic] = &kafka.Writer{
			Addr:         kafka.TCP(cfg.Brokers...),
			Topic:        topic,
			Balancer:     &kafka.LeastBytes{}, // 選擇負載最低的分區
			BatchSize:    100,                 // 批次大小
			BatchTimeout: 10 * time.Millisecond,
			RequiredAcks: kafka.RequireOne, // 等待 Leader 確認
			Async:        false,            // 同步發送，確保可靠性
		}
	}

	log.Info("kafka producer initialized",
		zap.Strings("brokers", cfg.Brokers),
		zap.Strings("topics", topics),
	)

	return &Producer{
		writers: writers,
		log:     log,
	}
}

// FlashSaleOrderMessage 秒殺訂單訊息（發送至 Kafka 供 Worker 處理）
type FlashSaleOrderMessage struct {
	MessageID   string    `json:"message_id"`
	Timestamp   time.Time `json:"timestamp"`
	FlashSaleID int64     `json:"flash_sale_id"`
	UserID      int64     `json:"user_id"`
	Quantity    int       `json:"quantity"`
	Ticket      string    `json:"ticket"` // 排隊憑證，用於前端查詢訂單狀態
}

// OrderStatusChangeMessage 訂單狀態變更訊息（用於 WebSocket 推送）
type OrderStatusChangeMessage struct {
	MessageID string    `json:"message_id"`
	Timestamp time.Time `json:"timestamp"`
	OrderNo   string    `json:"order_no"`
	UserID    int64     `json:"user_id"`
	OldStatus int       `json:"old_status"`
	NewStatus int       `json:"new_status"`
}

// AIAnalysisTaskMessage AI 分析任務訊息
type AIAnalysisTaskMessage struct {
	MessageID   string    `json:"message_id"`
	Timestamp   time.Time `json:"timestamp"`
	TaskType    string    `json:"task_type"` // strategy/anomaly/recommendation
	FlashSaleID int64     `json:"flash_sale_id,omitempty"`
	UserID      int64     `json:"user_id,omitempty"`
	Payload     string    `json:"payload,omitempty"`
}

// SendFlashSaleOrder 發送秒殺訂單訊息
func (p *Producer) SendFlashSaleOrder(ctx context.Context, msg *FlashSaleOrderMessage) error {
	return p.send(ctx, "flash-sale-orders", msg.UserID, msg)
}

// SendOrderStatusChange 發送訂單狀態變更訊息
func (p *Producer) SendOrderStatusChange(ctx context.Context, msg *OrderStatusChangeMessage) error {
	return p.send(ctx, "order-status-change", msg.UserID, msg)
}

// SendAIAnalysisTask 發送 AI 分析任務訊息
func (p *Producer) SendAIAnalysisTask(ctx context.Context, msg *AIAnalysisTaskMessage) error {
	return p.send(ctx, "ai-analysis-tasks", msg.FlashSaleID, msg)
}

// send 通用訊息發送方法
func (p *Producer) send(ctx context.Context, topic string, key int64, msg interface{}) error {
	writer, ok := p.writers[topic]
	if !ok {
		return fmt.Errorf("unknown topic: %s", topic)
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	err = writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(fmt.Sprintf("%d", key)), // 使用 UserID/FlashSaleID 作為分區 Key
		Value: data,
		Time:  time.Now(),
	})

	if err != nil {
		p.log.Error("failed to send kafka message",
			zap.String("topic", topic),
			zap.Error(err),
		)
		return err
	}

	p.log.Debug("kafka message sent",
		zap.String("topic", topic),
		zap.Int64("key", key),
	)

	return nil
}

// Close 關閉所有 Writer
func (p *Producer) Close() error {
	var errs []error
	for topic, writer := range p.writers {
		if err := writer.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close writer for %s: %w", topic, err))
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("errors closing producers: %v", errs)
	}
	return nil
}
