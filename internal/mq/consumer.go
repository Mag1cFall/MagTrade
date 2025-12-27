// Kafka 訊息消費者
//
// 本檔案封裝 Kafka 消費者功能，用於處理非同步任務
// 使用消費者群組（Consumer Group）保證訊息只被處理一次
// 支援註冊多個 Topic 的處理器，獨立 goroutine 消費
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

// MessageHandler 訊息處理函式類型
type MessageHandler func(ctx context.Context, msg []byte) error

// Consumer Kafka 消費者
type Consumer struct {
	readers  map[string]*kafka.Reader  // 每個 Topic 一個 Reader
	handlers map[string]MessageHandler // 每個 Topic 對應的處理函式
	log      *zap.Logger
	stopCh   chan struct{} // 停止信號
}

// NewConsumer 建立消費者
func NewConsumer(cfg *config.KafkaConfig, log *zap.Logger) *Consumer {
	readers := make(map[string]*kafka.Reader)

	topics := []string{
		cfg.Topics.FlashSaleOrders,
		cfg.Topics.OrderStatusChange,
		cfg.Topics.AIAnalysisTasks,
	}

	for _, topic := range topics {
		readers[topic] = kafka.NewReader(kafka.ReaderConfig{
			Brokers:        cfg.Brokers,
			Topic:          topic,
			GroupID:        cfg.ConsumerGroup, // 消費者群組，保證訊息不重複消費
			MinBytes:       10e3,              // 最小抓取量 10KB
			MaxBytes:       10e6,              // 最大抓取量 10MB
			MaxWait:        1 * time.Second,
			CommitInterval: 1 * time.Second,  // 自動提交 offset 間隔
			StartOffset:    kafka.LastOffset, // 從最新訊息開始消費
		})
	}

	log.Info("kafka consumer initialized",
		zap.Strings("brokers", cfg.Brokers),
		zap.String("group", cfg.ConsumerGroup),
		zap.Strings("topics", topics),
	)

	return &Consumer{
		readers:  readers,
		handlers: make(map[string]MessageHandler),
		log:      log,
		stopCh:   make(chan struct{}),
	}
}

// RegisterHandler 為指定 Topic 註冊處理函式
func (c *Consumer) RegisterHandler(topic string, handler MessageHandler) {
	c.handlers[topic] = handler
}

// Start 啟動所有消費者（每個 Topic 一個獨立 goroutine）
func (c *Consumer) Start(ctx context.Context) {
	for topic, reader := range c.readers {
		handler, ok := c.handlers[topic]
		if !ok {
			c.log.Warn("no handler registered for topic", zap.String("topic", topic))
			continue
		}

		go c.consume(ctx, topic, reader, handler)
	}
}

// consume 單一 Topic 的消費迴圈
func (c *Consumer) consume(ctx context.Context, topic string, reader *kafka.Reader, handler MessageHandler) {
	c.log.Info("starting consumer", zap.String("topic", topic))

	for {
		select {
		case <-ctx.Done(): // 收到關閉信號
			c.log.Info("stopping consumer", zap.String("topic", topic))
			return
		case <-c.stopCh:
			return
		default:
		}

		// 抓取下一條訊息
		msg, err := reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			c.log.Error("failed to fetch message",
				zap.String("topic", topic),
				zap.Error(err),
			)
			time.Sleep(time.Second) // 錯誤後短暫等待再重試
			continue
		}

		// 執行業務處理
		if err := handler(ctx, msg.Value); err != nil {
			c.log.Error("failed to handle message",
				zap.String("topic", topic),
				zap.String("key", string(msg.Key)),
				zap.Error(err),
			)
			continue // 處理失敗不提交 offset，下次會重試
		}

		// 提交 offset，標記訊息已處理
		if err := reader.CommitMessages(ctx, msg); err != nil {
			c.log.Error("failed to commit message",
				zap.String("topic", topic),
				zap.Error(err),
			)
		}

		c.log.Debug("message processed",
			zap.String("topic", topic),
			zap.String("key", string(msg.Key)),
		)
	}
}

// Stop 停止消費
func (c *Consumer) Stop() {
	close(c.stopCh)
}

// Close 關閉所有 Reader
func (c *Consumer) Close() error {
	c.Stop()
	var errs []error
	for topic, reader := range c.readers {
		if err := reader.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close reader for %s: %w", topic, err))
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("errors closing consumers: %v", errs)
	}
	return nil
}

// ParseFlashSaleOrderMessage 解析秒殺訂單訊息
func ParseFlashSaleOrderMessage(data []byte) (*FlashSaleOrderMessage, error) {
	var msg FlashSaleOrderMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

// ParseOrderStatusChangeMessage 解析訂單狀態變更訊息
func ParseOrderStatusChangeMessage(data []byte) (*OrderStatusChangeMessage, error) {
	var msg OrderStatusChangeMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

// ParseAIAnalysisTaskMessage 解析 AI 分析任務訊息
func ParseAIAnalysisTaskMessage(data []byte) (*AIAnalysisTaskMessage, error) {
	var msg AIAnalysisTaskMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}
