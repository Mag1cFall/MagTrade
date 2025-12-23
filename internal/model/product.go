package model

import (
	"time"

	"gorm.io/gorm"
)

type ProductStatus int8

const (
	ProductStatusOffShelf ProductStatus = 0
	ProductStatusOnShelf  ProductStatus = 1
)

type Product struct {
	ID            int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	Name          string         `gorm:"type:varchar(200);not null" json:"name"`
	Description   string         `gorm:"type:text" json:"description"`
	OriginalPrice float64        `gorm:"type:decimal(10,2);not null" json:"original_price"`
	ImageURL      string         `gorm:"type:text" json:"image_url"`
	Status        ProductStatus  `gorm:"type:smallint;default:1" json:"status"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Product) TableName() string {
	return "products"
}
