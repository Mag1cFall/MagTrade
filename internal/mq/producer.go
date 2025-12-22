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

type Producer struct {
	writers map[string]*kafka.Writer
	log     *zap.Logger
}

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
			Balancer:     &kafka.LeastBytes{},
			BatchSize:    100,
			BatchTimeout: 10 * time.Millisecond,
			RequiredAcks: kafka.RequireOne,
			Async:        false,
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

type FlashSaleOrderMessage struct {
	MessageID   string    `json:"message_id"`
	Timestamp   time.Time `json:"timestamp"`
	FlashSaleID int64     `json:"flash_sale_id"`
	UserID      int64     `json:"user_id"`
	Quantity    int       `json:"quantity"`
	Ticket      string    `json:"ticket"`
}

type OrderStatusChangeMessage struct {
	MessageID string    `json:"message_id"`
	Timestamp time.Time `json:"timestamp"`
	OrderNo   string    `json:"order_no"`
	UserID    int64     `json:"user_id"`
	OldStatus int       `json:"old_status"`
	NewStatus int       `json:"new_status"`
}

type AIAnalysisTaskMessage struct {
	MessageID   string    `json:"message_id"`
	Timestamp   time.Time `json:"timestamp"`
	TaskType    string    `json:"task_type"`
	FlashSaleID int64     `json:"flash_sale_id,omitempty"`
	UserID      int64     `json:"user_id,omitempty"`
	Payload     string    `json:"payload,omitempty"`
}

func (p *Producer) SendFlashSaleOrder(ctx context.Context, msg *FlashSaleOrderMessage) error {
	return p.send(ctx, "flash-sale-orders", msg.UserID, msg)
}

func (p *Producer) SendOrderStatusChange(ctx context.Context, msg *OrderStatusChangeMessage) error {
	return p.send(ctx, "order-status-change", msg.UserID, msg)
}

func (p *Producer) SendAIAnalysisTask(ctx context.Context, msg *AIAnalysisTaskMessage) error {
	return p.send(ctx, "ai-analysis-tasks", msg.FlashSaleID, msg)
}

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
		Key:   []byte(fmt.Sprintf("%d", key)),
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
