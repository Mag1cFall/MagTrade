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
	ErrFlashSaleNotFound = errors.New("flash sale not found")
)

type FlashSaleRepository struct {
	db *gorm.DB
}

func NewFlashSaleRepository() *FlashSaleRepository {
	return &FlashSaleRepository{db: database.Get()}
}

func (r *FlashSaleRepository) Create(ctx context.Context, flashSale *model.FlashSale) error {
	return r.db.WithContext(ctx).Create(flashSale).Error
}

func (r *FlashSaleRepository) GetByID(ctx context.Context, id int64) (*model.FlashSale, error) {
	var flashSale model.FlashSale
	result := r.db.WithContext(ctx).Preload("Product").First(&flashSale, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrFlashSaleNotFound
		}
		return nil, result.Error
	}
	return &flashSale, nil
}

func (r *FlashSaleRepository) List(ctx context.Context, page, pageSize int, status *model.FlashSaleStatus) ([]model.FlashSale, int64, error) {
	var flashSales []model.FlashSale
	var total int64

	db := r.db.WithContext(ctx).Model(&model.FlashSale{})

	if status != nil {
		db = db.Where("status = ?", *status)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := db.Preload("Product").Offset(offset).Limit(pageSize).Order("start_time DESC").Find(&flashSales).Error; err != nil {
		return nil, 0, err
	}

	return flashSales, total, nil
}

func (r *FlashSaleRepository) ListActive(ctx context.Context) ([]model.FlashSale, error) {
	var flashSales []model.FlashSale
	now := time.Now()

	result := r.db.WithContext(ctx).
		Preload("Product").
		Where("status = ? AND start_time <= ? AND end_time > ?",
			model.FlashSaleStatusActive, now, now).
		Order("start_time ASC").
		Find(&flashSales)

	return flashSales, result.Error
}

func (r *FlashSaleRepository) ListUpcoming(ctx context.Context, limit int) ([]model.FlashSale, error) {
	var flashSales []model.FlashSale
	now := time.Now()

	result := r.db.WithContext(ctx).
		Preload("Product").
		Where("status = ? AND start_time > ?", model.FlashSaleStatusPending, now).
		Order("start_time ASC").
		Limit(limit).
		Find(&flashSales)

	return flashSales, result.Error
}

func (r *FlashSaleRepository) Update(ctx context.Context, flashSale *model.FlashSale) error {
	return r.db.WithContext(ctx).Save(flashSale).Error
}

func (r *FlashSaleRepository) DecrementStock(ctx context.Context, id int64, quantity int) error {
	result := r.db.WithContext(ctx).
		Model(&model.FlashSale{}).
		Where("id = ? AND available_stock >= ?", id, quantity).
		UpdateColumn("available_stock", gorm.Expr("available_stock - ?", quantity))

	if result.RowsAffected == 0 {
		return errors.New("insufficient stock or flash sale not found")
	}

	return result.Error
}

func (r *FlashSaleRepository) IncrementStock(ctx context.Context, id int64, quantity int) error {
	return r.db.WithContext(ctx).
		Model(&model.FlashSale{}).
		Where("id = ?", id).
		UpdateColumn("available_stock", gorm.Expr("available_stock + ?", quantity)).
		Error
}

func (r *FlashSaleRepository) UpdatePendingToActive(ctx context.Context) (int64, error) {
	now := time.Now()
	result := r.db.WithContext(ctx).
		Model(&model.FlashSale{}).
		Where("status = ? AND start_time <= ?", model.FlashSaleStatusPending, now).
		Update("status", model.FlashSaleStatusActive)

	return result.RowsAffected, result.Error
}

func (r *FlashSaleRepository) UpdateActiveToFinished(ctx context.Context) (int64, error) {
	now := time.Now()
	result := r.db.WithContext(ctx).
		Model(&model.FlashSale{}).
		Where("status = ? AND end_time <= ?", model.FlashSaleStatusActive, now).
		Update("status", model.FlashSaleStatusFinished)

	return result.RowsAffected, result.Error
}
