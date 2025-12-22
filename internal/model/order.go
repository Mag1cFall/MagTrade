package model

import (
	"time"

	"gorm.io/gorm"
)

type OrderStatus int8

const (
	OrderStatusPending   OrderStatus = 0
	OrderStatusPaid      OrderStatus = 1
	OrderStatusCancelled OrderStatus = 2
	OrderStatusRefunded  OrderStatus = 3
)

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

type Order struct {
	ID          int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderNo     string         `gorm:"type:varchar(32);uniqueIndex;not null" json:"order_no"`
	UserID      int64          `gorm:"index;not null" json:"user_id"`
	User        *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	FlashSaleID int64          `gorm:"index;not null" json:"flash_sale_id"`
	FlashSale   *FlashSale     `gorm:"foreignKey:FlashSaleID" json:"flash_sale,omitempty"`
	Amount      float64        `gorm:"type:decimal(10,2);not null" json:"amount"`
	Quantity    int            `gorm:"default:1" json:"quantity"`
	Status      OrderStatus    `gorm:"type:smallint;default:0" json:"status"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	PaidAt      *time.Time     `json:"paid_at,omitempty"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Order) TableName() string {
	return "orders"
}

func (o *Order) CanPay() bool {
	return o.Status == OrderStatusPending
}

func (o *Order) CanCancel() bool {
	return o.Status == OrderStatusPending
}

func (o *Order) CanRefund() bool {
	return o.Status == OrderStatusPaid
}
