package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/Mag1cFall/magtrade/internal/config"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var rdb *redis.Client

func Init(cfg *config.RedisConfig, log *zap.Logger) error {
	rdb = redis.NewClient(&redis.Options{
		Addr:         cfg.Addr(),
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: 10,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect redis: %w", err)
	}

	log.Info("redis connected",
		zap.String("addr", cfg.Addr()),
		zap.Int("db", cfg.DB),
	)

	return nil
}

func Get() *redis.Client {
	if rdb == nil {
		panic("redis not initialized, call Init() first")
	}
	return rdb
}

func Close() error {
	if rdb == nil {
		return nil
	}
	return rdb.Close()
}
