// 審計日誌資料模型
//
// 對應資料表 audit_logs，記錄使用者關鍵操作軌跡
// 用於安全審計、問題追蹤和合規需求
package model

import (
	"time"

	"gorm.io/gorm"
)

// AuditLog 審計日誌模型
type AuditLog struct {
	ID         int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID     int64          `gorm:"index" json:"user_id"` // 可為空（匿名請求）
	Action     string         `gorm:"type:varchar(50);not null" json:"action"`
	Resource   string         `gorm:"type:varchar(100)" json:"resource"` // 資源類型：user/order/flash_sale
	ResourceID string         `gorm:"type:varchar(50)" json:"resource_id"`
	IP         string         `gorm:"type:varchar(45)" json:"ip"` // 支援 IPv6
	UserAgent  string         `gorm:"type:varchar(255)" json:"user_agent"`
	Details    string         `gorm:"type:text" json:"details"` // JSON 格式額外資訊
	CreatedAt  time.Time      `gorm:"autoCreateTime" json:"created_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}

// 預定義的審計動作常量
const (
	AuditActionLogin          = "login"           // 登入
	AuditActionLogout         = "logout"          // 登出
	AuditActionRegister       = "register"        // 註冊
	AuditActionPasswordChange = "password_change" // 修改密碼
	AuditActionOrderCreate    = "order_create"    // 建立訂單
	AuditActionOrderPay       = "order_pay"       // 支付訂單
	AuditActionOrderCancel    = "order_cancel"    // 取消訂單
	AuditActionFlashSaleRush  = "flash_sale_rush" // 秒殺搶購
)
