package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type StockService struct {
	rdb *redis.Client
}

func NewStockService() *StockService {
	return &StockService{rdb: Get()}
}

func StockKey(flashSaleID int64) string {
	return fmt.Sprintf("flash:stock:%d", flashSaleID)
}

func BoughtKey(flashSaleID, userID int64) string {
	return fmt.Sprintf("flash:bought:%d:%d", flashSaleID, userID)
}

func LockKey(flashSaleID, userID int64) string {
	return fmt.Sprintf("flash:lock:%d:%d", flashSaleID, userID)
}

func (s *StockService) InitStock(ctx context.Context, flashSaleID int64, stock int) error {
	key := StockKey(flashSaleID)
	return s.rdb.Set(ctx, key, stock, 24*time.Hour).Err()
}

func (s *StockService) GetStock(ctx context.Context, flashSaleID int64) (int, error) {
	key := StockKey(flashSaleID)
	val, err := s.rdb.Get(ctx, key).Int()
	if err == redis.Nil {
		return 0, nil
	}
	return val, err
}

type DeductResult struct {
	Success bool
	Code    int
	Message string
}

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

func (s *StockService) RestoreStock(ctx context.Context, flashSaleID, userID int64, quantity int) error {
	stockKey := StockKey(flashSaleID)
	boughtKey := BoughtKey(flashSaleID, userID)

	_, err := s.rdb.Eval(ctx, RestoreStockScript,
		[]string{stockKey, boughtKey},
		quantity,
	).Result()

	return err
}

type DistributedLock struct {
	rdb   *redis.Client
	key   string
	value string
}

func NewDistributedLock(flashSaleID, userID int64) *DistributedLock {
	return &DistributedLock{
		rdb:   Get(),
		key:   LockKey(flashSaleID, userID),
		value: uuid.New().String(),
	}
}

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

func (l *DistributedLock) Unlock(ctx context.Context) error {
	_, err := l.rdb.Eval(ctx, DistributedUnlockScript,
		[]string{l.key},
		l.value,
	).Result()

	return err
}
