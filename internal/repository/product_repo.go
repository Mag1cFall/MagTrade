// 商品資料存取層
//
// 本檔案封裝商品表的 CRUD 操作
// 列表查詢只返回上架狀態的商品
package repository

import (
	"context"
	"errors"

	"github.com/Mag1cFall/magtrade/internal/database"
	"github.com/Mag1cFall/magtrade/internal/model"
	"gorm.io/gorm"
)

var (
	ErrProductNotFound = errors.New("product not found")
)

// ProductRepository 商品資料存取
type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository() *ProductRepository {
	return &ProductRepository{db: database.Get()}
}

// Create 建立商品
func (r *ProductRepository) Create(ctx context.Context, product *model.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

// GetByID 根據 ID 查詢商品
func (r *ProductRepository) GetByID(ctx context.Context, id int64) (*model.Product, error) {
	var product model.Product
	result := r.db.WithContext(ctx).First(&product, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, result.Error
	}
	return &product, nil
}

// List 分頁查詢商品（只返回上架狀態）
func (r *ProductRepository) List(ctx context.Context, page, pageSize int) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64

	db := r.db.WithContext(ctx).Model(&model.Product{}).Where("status = ?", model.ProductStatusOnShelf)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := db.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// Update 更新商品
func (r *ProductRepository) Update(ctx context.Context, product *model.Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}

// Delete 刪除商品（GORM 軟刪除）
func (r *ProductRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.Product{}, id).Error
}
