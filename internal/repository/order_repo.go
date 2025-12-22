package repository

import (
	"context"
	"errors"
	"time"

	"github.com/Mag1cFall/magtrade/internal/database"
	"github.com/Mag1cFall/magtrade/internal/model"
	"gorm.io/gorm"
)

var (
	ErrOrderNotFound = errors.New("order not found")
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{db: database.Get()}
}

func (r *OrderRepository) Create(ctx context.Context, order *model.Order) error {
	return r.db.WithContext(ctx).Create(order).Error
}

func (r *OrderRepository) GetByID(ctx context.Context, id int64) (*model.Order, error) {
	var order model.Order
	result := r.db.WithContext(ctx).Preload("FlashSale.Product").First(&order, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrOrderNotFound
		}
		return nil, result.Error
	}
	return &order, nil
}

func (r *OrderRepository) GetByOrderNo(ctx context.Context, orderNo string) (*model.Order, error) {
	var order model.Order
	result := r.db.WithContext(ctx).Preload("FlashSale.Product").Where("order_no = ?", orderNo).First(&order)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrOrderNotFound
		}
		return nil, result.Error
	}
	return &order, nil
}

func (r *OrderRepository) GetByUserAndFlashSale(ctx context.Context, userID, flashSaleID int64) (*model.Order, error) {
	var order model.Order
	result := r.db.WithContext(ctx).
		Where("user_id = ? AND flash_sale_id = ? AND status != ?", userID, flashSaleID, model.OrderStatusCancelled).
		First(&order)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &order, nil
}

func (r *OrderRepository) ListByUser(ctx context.Context, userID int64, page, pageSize int) ([]model.Order, int64, error) {
	var orders []model.Order
	var total int64

	db := r.db.WithContext(ctx).Model(&model.Order{}).Where("user_id = ?", userID)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := db.Preload("FlashSale.Product").Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (r *OrderRepository) UpdateStatus(ctx context.Context, id int64, oldStatus, newStatus model.OrderStatus) error {
	result := r.db.WithContext(ctx).
		Model(&model.Order{}).
		Where("id = ? AND status = ?", id, oldStatus).
		Updates(map[string]interface{}{
			"status":     newStatus,
			"updated_at": time.Now(),
		})

	if result.RowsAffected == 0 {
		return errors.New("order status update failed: status mismatch or order not found")
	}

	return result.Error
}

func (r *OrderRepository) Pay(ctx context.Context, id int64) error {
	now := time.Now()
	result := r.db.WithContext(ctx).
		Model(&model.Order{}).
		Where("id = ? AND status = ?", id, model.OrderStatusPending).
		Updates(map[string]interface{}{
			"status":     model.OrderStatusPaid,
			"paid_at":    now,
			"updated_at": now,
		})

	if result.RowsAffected == 0 {
		return errors.New("order pay failed: not pending or not found")
	}

	return result.Error
}

func (r *OrderRepository) Cancel(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).
		Model(&model.Order{}).
		Where("id = ? AND status = ?", id, model.OrderStatusPending).
		Updates(map[string]interface{}{
			"status":     model.OrderStatusCancelled,
			"updated_at": time.Now(),
		})

	if result.RowsAffected == 0 {
		return errors.New("order cancel failed: not pending or not found")
	}

	return result.Error
}

func (r *OrderRepository) CountExpiredPending(ctx context.Context, expireDuration time.Duration) (int64, error) {
	var count int64
	expireTime := time.Now().Add(-expireDuration)

	result := r.db.WithContext(ctx).
		Model(&model.Order{}).
		Where("status = ? AND created_at < ?", model.OrderStatusPending, expireTime).
		Count(&count)

	return count, result.Error
}

func (r *OrderRepository) CancelExpiredPending(ctx context.Context, expireDuration time.Duration, limit int) ([]model.Order, error) {
	var orders []model.Order
	expireTime := time.Now().Add(-expireDuration)

	if err := r.db.WithContext(ctx).
		Where("status = ? AND created_at < ?", model.OrderStatusPending, expireTime).
		Limit(limit).
		Find(&orders).Error; err != nil {
		return nil, err
	}

	if len(orders) == 0 {
		return orders, nil
	}

	var ids []int64
	for _, o := range orders {
		ids = append(ids, o.ID)
	}

	if err := r.db.WithContext(ctx).
		Model(&model.Order{}).
		Where("id IN ?", ids).
		Updates(map[string]interface{}{
			"status":     model.OrderStatusCancelled,
			"updated_at": time.Now(),
		}).Error; err != nil {
		return nil, err
	}

	return orders, nil
}
