package model

import (
	"time"

	"gorm.io/gorm"
)

type AuditLog struct {
	ID         int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID     int64          `gorm:"index" json:"user_id"`
	Action     string         `gorm:"type:varchar(50);not null" json:"action"`
	Resource   string         `gorm:"type:varchar(100)" json:"resource"`
	ResourceID string         `gorm:"type:varchar(50)" json:"resource_id"`
	IP         string         `gorm:"type:varchar(45)" json:"ip"`
	UserAgent  string         `gorm:"type:varchar(255)" json:"user_agent"`
	Details    string         `gorm:"type:text" json:"details"`
	CreatedAt  time.Time      `gorm:"autoCreateTime" json:"created_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}

const (
	AuditActionLogin          = "login"
	AuditActionLogout         = "logout"
	AuditActionRegister       = "register"
	AuditActionPasswordChange = "password_change"
	AuditActionOrderCreate    = "order_create"
	AuditActionOrderPay       = "order_pay"
	AuditActionOrderCancel    = "order_cancel"
	AuditActionFlashSaleRush  = "flash_sale_rush"
)
