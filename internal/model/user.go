// 使用者資料模型
//
// 對應資料表 users，支援 GORM 軟刪除
// 角色分為 user（普通使用者）和 admin（管理員）
package model

import (
	"time"

	"gorm.io/gorm"
)

// UserStatus 使用者狀態
type UserStatus int8

const (
	UserStatusDisabled UserStatus = 0 // 停用
	UserStatusActive   UserStatus = 1 // 啟用
)

// UserRole 使用者角色
type UserRole string

const (
	UserRoleUser  UserRole = "user"
	UserRoleAdmin UserRole = "admin"
)

// User 使用者模型
type User struct {
	ID            int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	Username      string         `gorm:"type:varchar(50);uniqueIndex:idx_users_username;not null" json:"username"`
	Email         string         `gorm:"type:varchar(100);uniqueIndex:idx_users_email;not null" json:"email"`
	PasswordHash  string         `gorm:"type:varchar(255);not null" json:"-"` // json:"-" 防止序列化輸出
	Role          UserRole       `gorm:"type:varchar(20);default:user" json:"role"`
	Status        UserStatus     `gorm:"type:smallint;default:1" json:"status"`
	EmailVerified bool           `gorm:"default:false" json:"email_verified"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index:idx_users_deleted_at" json:"-"` // GORM 軟刪除
}

// TableName 指定資料表名稱
func (User) TableName() string {
	return "users"
}
