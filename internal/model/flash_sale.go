package model

import (
	"time"

	"gorm.io/gorm"
)

type FlashSaleStatus int8

const (
	FlashSaleStatusPending  FlashSaleStatus = 0
	FlashSaleStatusActive   FlashSaleStatus = 1
	FlashSaleStatusFinished FlashSaleStatus = 2
)

type FlashSale struct {
	ID             int64           `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductID      int64           `gorm:"index;not null" json:"product_id"`
	Product        *Product        `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	FlashPrice     float64         `gorm:"type:decimal(10,2);not null" json:"flash_price"`
	TotalStock     int             `gorm:"not null" json:"total_stock"`
	AvailableStock int             `gorm:"not null" json:"available_stock"`
	PerUserLimit   int             `gorm:"default:1" json:"per_user_limit"`
	StartTime      time.Time       `gorm:"not null;index" json:"start_time"`
	EndTime        time.Time       `gorm:"not null;index" json:"end_time"`
	Status         FlashSaleStatus `gorm:"type:smallint;default:0" json:"status"`
	CreatedAt      time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt      gorm.DeletedAt  `gorm:"index" json:"-"`
}

func (FlashSale) TableName() string {
	return "flash_sales"
}

func (f *FlashSale) IsActive() bool {
	now := time.Now()
	return f.Status == FlashSaleStatusActive &&
		now.After(f.StartTime) &&
		now.Before(f.EndTime)
}

func (f *FlashSale) IsPending() bool {
	now := time.Now()
	return f.Status == FlashSaleStatusPending && now.Before(f.StartTime)
}
