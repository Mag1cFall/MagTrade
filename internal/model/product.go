// 商品資料模型
//
// 對應資料表 products，儲存可銷售的商品資訊
// 商品可關聯至秒殺活動（FlashSale）
package model

import (
	"time"

	"gorm.io/gorm"
)

// ProductStatus 商品狀態
type ProductStatus int8

const (
	ProductStatusOffShelf ProductStatus = 0 // 下架
	ProductStatusOnShelf  ProductStatus = 1 // 上架
)

// Product 商品模型
type Product struct {
	ID            int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	Name          string         `gorm:"type:varchar(200);not null" json:"name"`
	Description   string         `gorm:"type:text" json:"description"`
	OriginalPrice float64        `gorm:"type:decimal(10,2);not null" json:"original_price"` // 單位：元
	ImageURL      string         `gorm:"type:text" json:"image_url"`
	Status        ProductStatus  `gorm:"type:smallint;default:1" json:"status"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定資料表名稱
func (Product) TableName() string {
	return "products"
}
