// Redis 客戶端初始化與連線管理
//
// 本檔案提供 Redis 連線池的初始化、獲取和關閉功能
// 使用 go-redis/v9 客戶端，支援連線池和超時配置
package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/Mag1cFall/magtrade/internal/config"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var rdb *redis.Client // 全域 Redis 客戶端實例

// Init 初始化 Redis 連線池
func Init(cfg *config.RedisConfig, log *zap.Logger) error {
	rdb = redis.NewClient(&redis.Options{
		Addr:         cfg.Addr(),
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize, // 連線池大小，生產環境建議 100-200
		MinIdleConns: 10,           // 最小閒置連線數，避免冷啟動延遲
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	// 測試連線是否正常
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

// Get 取得 Redis 客戶端實例，若未初始化則 panic
func Get() *redis.Client {
	if rdb == nil {
		panic("redis not initialized, call Init() first")
	}
	return rdb
}

// Close 關閉 Redis 連線
func Close() error {
	if rdb == nil {
		return nil
	}
	return rdb.Close()
}

// GetClient 取得客戶端（Get 的別名，兼容舊代碼）
func GetClient() *redis.Client {
	return Get()
}
