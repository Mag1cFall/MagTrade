// 庫存快取服務與分散式鎖
//
// 本檔案提供秒殺核心的快取操作功能
// StockService: 庫存的初始化、查詢、扣減、恢復
// DistributedLock: 基於 Redis 的分散式鎖，防止重複提交
package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// StockService 庫存快取服務
type StockService struct {
	rdb *redis.Client
}

func NewStockService() *StockService {
	return &StockService{rdb: Get()}
}

// StockKey 生成庫存 Redis Key，格式: flash:stock:{活動ID}
func StockKey(flashSaleID int64) string {
	return fmt.Sprintf("flash:stock:%d", flashSaleID)
}

// BoughtKey 生成使用者已購數量 Key，格式: flash:bought:{活動ID}:{使用者ID}
func BoughtKey(flashSaleID, userID int64) string {
	return fmt.Sprintf("flash:bought:%d:%d", flashSaleID, userID)
}

// LockKey 生成分散式鎖 Key，格式: flash:lock:{活動ID}:{使用者ID}
func LockKey(flashSaleID, userID int64) string {
	return fmt.Sprintf("flash:lock:%d:%d", flashSaleID, userID)
}

// InitStock 初始化秒殺活動庫存至 Redis，有效期 24 小時
func (s *StockService) InitStock(ctx context.Context, flashSaleID int64, stock int) error {
	key := StockKey(flashSaleID)
	return s.rdb.Set(ctx, key, stock, 24*time.Hour).Err()
}

// GetStock 查詢當前庫存數量
func (s *StockService) GetStock(ctx context.Context, flashSaleID int64) (int, error) {
	key := StockKey(flashSaleID)
	val, err := s.rdb.Get(ctx, key).Int()
	if err == redis.Nil { // Key 不存在視為庫存 0
		return 0, nil
	}
	return val, err
}

// DeductResult 庫存扣減結果
type DeductResult struct {
	Success bool
	Code    int // 1=成功, -1=庫存不足, -2=超過限購
	Message string
}

// DeductStock 扣減庫存（原子操作）
// 使用 Lua 腳本在 Redis 端執行，保證併發安全
func (s *StockService) DeductStock(ctx context.Context, flashSaleID, userID int64, quantity, limit int) (*DeductResult, error) {
	stockKey := StockKey(flashSaleID)
	boughtKey := BoughtKey(flashSaleID, userID)

	result, err := s.rdb.Eval(ctx, DeductStockScript,
		[]string{stockKey, boughtKey},
		quantity, limit,
	).Result()

	if err != nil {
		return nil, fmt.Errorf("failed to execute deduct script: %w", err)
	}

	// 解析 Lua 腳本返回值 [code, message]
	arr, ok := result.([]interface{})
	if !ok || len(arr) != 2 {
		return nil, fmt.Errorf("unexpected script result format")
	}

	code, _ := arr[0].(int64)
	msg, _ := arr[1].(string)

	return &DeductResult{
		Success: code == 1,
		Code:    int(code),
		Message: msg,
	}, nil
}

// RestoreStock 恢復庫存（訂單取消或發送失敗時使用）
func (s *StockService) RestoreStock(ctx context.Context, flashSaleID, userID int64, quantity int) error {
	stockKey := StockKey(flashSaleID)
	boughtKey := BoughtKey(flashSaleID, userID)

	_, err := s.rdb.Eval(ctx, RestoreStockScript,
		[]string{stockKey, boughtKey},
		quantity,
	).Result()

	return err
}

// DistributedLock 分散式鎖，防止同一使用者重複提交
type DistributedLock struct {
	rdb   *redis.Client
	key   string
	value string // UUID，確保只有自己能解鎖
}

// NewDistributedLock 建立分散式鎖實例
func NewDistributedLock(flashSaleID, userID int64) *DistributedLock {
	return &DistributedLock{
		rdb:   Get(),
		key:   LockKey(flashSaleID, userID),
		value: uuid.New().String(), // 每次鎖用不同 UUID
	}
}

// Lock 嘗試加鎖，expireMs 為鎖的過期時間（毫秒）
// 返回 true 表示加鎖成功，false 表示鎖已被其他人持有
func (l *DistributedLock) Lock(ctx context.Context, expireMs int) (bool, error) {
	result, err := l.rdb.Eval(ctx, DistributedLockScript,
		[]string{l.key},
		l.value, expireMs,
	).Result()

	if err != nil {
		return false, err
	}

	return result.(int64) == 1, nil
}

// Unlock 釋放鎖（只有鎖的持有者才能成功釋放）
func (l *DistributedLock) Unlock(ctx context.Context) error {
	_, err := l.rdb.Eval(ctx, DistributedUnlockScript,
		[]string{l.key},
		l.value,
	).Result()

	return err
}
