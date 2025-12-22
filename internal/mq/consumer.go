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

type MessageHandler func(ctx context.Context, msg []byte) error

type Consumer struct {
	readers  map[string]*kafka.Reader
	handlers map[string]MessageHandler
	log      *zap.Logger
	stopCh   chan struct{}
}

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
			GroupID:        cfg.ConsumerGroup,
			MinBytes:       10e3,
			MaxBytes:       10e6,
			MaxWait:        1 * time.Second,
			CommitInterval: 1 * time.Second,
			StartOffset:    kafka.LastOffset,
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

func (c *Consumer) RegisterHandler(topic string, handler MessageHandler) {
	c.handlers[topic] = handler
}

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

func (c *Consumer) consume(ctx context.Context, topic string, reader *kafka.Reader, handler MessageHandler) {
	c.log.Info("starting consumer", zap.String("topic", topic))

	for {
		select {
		case <-ctx.Done():
			c.log.Info("stopping consumer", zap.String("topic", topic))
			return
		case <-c.stopCh:
			return
		default:
		}

		msg, err := reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			c.log.Error("failed to fetch message",
				zap.String("topic", topic),
				zap.Error(err),
			)
			time.Sleep(time.Second)
			continue
		}

		if err := handler(ctx, msg.Value); err != nil {
			c.log.Error("failed to handle message",
				zap.String("topic", topic),
				zap.String("key", string(msg.Key)),
				zap.Error(err),
			)
			continue
		}

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

func (c *Consumer) Stop() {
	close(c.stopCh)
}

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

func ParseFlashSaleOrderMessage(data []byte) (*FlashSaleOrderMessage, error) {
	var msg FlashSaleOrderMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

func ParseOrderStatusChangeMessage(data []byte) (*OrderStatusChangeMessage, error) {
	var msg OrderStatusChangeMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

func ParseAIAnalysisTaskMessage(data []byte) (*AIAnalysisTaskMessage, error) {
	var msg AIAnalysisTaskMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}
