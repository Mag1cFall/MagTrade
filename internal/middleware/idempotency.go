// 冪等性檢查器
//
// 防止使用者重複提交相同請求（如重複下單）
// 使用 Redis SETNX 實現分散式冪等性檢查
package middleware

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// IdempotencyChecker 冪等性檢查器
type IdempotencyChecker struct {
	rdb *redis.Client
	ttl time.Duration // 冪等 Key 過期時間
}

func NewIdempotencyChecker(rdb *redis.Client) *IdempotencyChecker {
	return &IdempotencyChecker{
		rdb: rdb,
		ttl: 10 * time.Minute, // 10 分鐘內相同請求視為重複
	}
}

// GenerateKey 根據使用者、動作、資源生成唯一冪等 Key
func (c *IdempotencyChecker) GenerateKey(userID int64, action string, resourceID int64) string {
	raw := fmt.Sprintf("%d:%s:%d", userID, action, resourceID)
	hash := sha256.Sum256([]byte(raw)) // SHA256 縮短 Key 長度
	return "idempotent:" + hex.EncodeToString(hash[:])
}

// Check 檢查 Key 是否存在，不存在則標記
// 返回 true 表示首次請求，false 表示重複請求
func (c *IdempotencyChecker) Check(ctx context.Context, key string) (bool, error) {
	result, err := c.rdb.SetNX(ctx, key, "1", c.ttl).Result() // SetNX = SET if Not eXists
	if err != nil {
		return false, err
	}
	return result, nil
}

// CheckAndMark 生成 Key 並檢查冪等性
func (c *IdempotencyChecker) CheckAndMark(ctx context.Context, userID int64, action string, resourceID int64) (bool, error) {
	key := c.GenerateKey(userID, action, resourceID)
	return c.Check(ctx, key)
}

// Remove 刪除冪等 Key（用於業務失敗時回滾）
func (c *IdempotencyChecker) Remove(ctx context.Context, key string) error {
	return c.rdb.Del(ctx, key).Err()
}
