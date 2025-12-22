package middleware

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type IdempotencyChecker struct {
	rdb *redis.Client
	ttl time.Duration
}

func NewIdempotencyChecker(rdb *redis.Client) *IdempotencyChecker {
	return &IdempotencyChecker{
		rdb: rdb,
		ttl: 10 * time.Minute,
	}
}

func (c *IdempotencyChecker) GenerateKey(userID int64, action string, resourceID int64) string {
	raw := fmt.Sprintf("%d:%s:%d", userID, action, resourceID)
	hash := sha256.Sum256([]byte(raw))
	return "idempotent:" + hex.EncodeToString(hash[:])
}

func (c *IdempotencyChecker) Check(ctx context.Context, key string) (bool, error) {
	result, err := c.rdb.SetNX(ctx, key, "1", c.ttl).Result()
	if err != nil {
		return false, err
	}
	return result, nil
}

func (c *IdempotencyChecker) CheckAndMark(ctx context.Context, userID int64, action string, resourceID int64) (bool, error) {
	key := c.GenerateKey(userID, action, resourceID)
	return c.Check(ctx, key)
}

func (c *IdempotencyChecker) Remove(ctx context.Context, key string) error {
	return c.rdb.Del(ctx, key).Err()
}
