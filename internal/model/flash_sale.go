// 秒殺活動資料模型
//
// 對應資料表 flash_sales，儲存限時搶購活動資訊
// 狀態流轉：待開始(0) → 進行中(1) → 已結束(2)
// 庫存由 Redis 和 DB 雙寫，以 Redis 為準
package model

import (
	"time"

	"gorm.io/gorm"
)

// FlashSaleStatus 秒殺活動狀態
type FlashSaleStatus int8

const (
	FlashSaleStatusPending  FlashSaleStatus = 0 // 待開始
	FlashSaleStatusActive   FlashSaleStatus = 1 // 進行中
	FlashSaleStatusFinished FlashSaleStatus = 2 // 已結束
)

// FlashSale 秒殺活動模型
type FlashSale struct {
	ID             int64           `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductID      int64           `gorm:"index;not null" json:"product_id"`
	Product        *Product        `gorm:"foreignKey:ProductID" json:"product,omitempty"` // GORM 關聯
	FlashPrice     float64         `gorm:"type:decimal(10,2);not null" json:"flash_price"`
	TotalStock     int             `gorm:"not null" json:"total_stock"`     // 總庫存
	AvailableStock int             `gorm:"not null" json:"available_stock"` // 剩餘庫存（DB）
	PerUserLimit   int             `gorm:"default:1" json:"per_user_limit"` // 每人限購數量
	StartTime      time.Time       `gorm:"not null;index" json:"start_time"`
	EndTime        time.Time       `gorm:"not null;index" json:"end_time"`
	Status         FlashSaleStatus `gorm:"type:smallint;default:0" json:"status"`
	CreatedAt      time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt      gorm.DeletedAt  `gorm:"index" json:"-"`
}

// TableName 指定資料表名稱
func (FlashSale) TableName() string {
	return "flash_sales"
}

// IsActive 檢查活動是否正在進行中
func (f *FlashSale) IsActive() bool {
	now := time.Now()
	return f.Status == FlashSaleStatusActive &&
		now.After(f.StartTime) &&
		now.Before(f.EndTime)
}

// IsPending 檢查活動是否待開始
func (f *FlashSale) IsPending() bool {
	now := time.Now()
	return f.Status == FlashSaleStatusPending && now.Before(f.StartTime)
}
