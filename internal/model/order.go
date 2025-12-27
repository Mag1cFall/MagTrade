// 訂單資料模型
//
// 對應資料表 orders，儲存使用者秒殺訂單
// 狀態流轉：待付款(0) → 已付款(1) / 已取消(2)
// 已付款訂單可退款 → 已退款(3)
package model

import (
	"time"

	"gorm.io/gorm"
)

// OrderStatus 訂單狀態
type OrderStatus int8

const (
	OrderStatusPending   OrderStatus = 0 // 待付款
	OrderStatusPaid      OrderStatus = 1 // 已付款
	OrderStatusCancelled OrderStatus = 2 // 已取消
	OrderStatusRefunded  OrderStatus = 3 // 已退款
)

// String 返回狀態的字串表示
func (s OrderStatus) String() string {
	switch s {
	case OrderStatusPending:
		return "pending"
	case OrderStatusPaid:
		return "paid"
	case OrderStatusCancelled:
		return "cancelled"
	case OrderStatusRefunded:
		return "refunded"
	default:
		return "unknown"
	}
}

// Order 訂單模型
type Order struct {
	ID          int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderNo     string         `gorm:"type:varchar(32);uniqueIndex;not null" json:"order_no"` // 雪花演算法生成
	UserID      int64          `gorm:"index;not null" json:"user_id"`
	User        *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	FlashSaleID int64          `gorm:"index;not null" json:"flash_sale_id"`
	FlashSale   *FlashSale     `gorm:"foreignKey:FlashSaleID" json:"flash_sale,omitempty"`
	Amount      float64        `gorm:"type:decimal(10,2);not null" json:"amount"` // 訂單金額
	Quantity    int            `gorm:"default:1" json:"quantity"`
	Status      OrderStatus    `gorm:"type:smallint;default:0" json:"status"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	PaidAt      *time.Time     `json:"paid_at,omitempty"` // 付款時間
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定資料表名稱
func (Order) TableName() string {
	return "orders"
}

// CanPay 檢查訂單是否可付款
func (o *Order) CanPay() bool {
	return o.Status == OrderStatusPending
}

// CanCancel 檢查訂單是否可取消
func (o *Order) CanCancel() bool {
	return o.Status == OrderStatusPending
}

// CanRefund 檢查訂單是否可退款
func (o *Order) CanRefund() bool {
	return o.Status == OrderStatusPaid
}
