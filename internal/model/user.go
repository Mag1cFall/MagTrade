package model

import (
	"time"

	"gorm.io/gorm"
)

type UserStatus int8

const (
	UserStatusDisabled UserStatus = 0
	UserStatusActive   UserStatus = 1
)

type UserRole string

const (
	UserRoleUser  UserRole = "user"
	UserRoleAdmin UserRole = "admin"
)

type User struct {
	ID            int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	Username      string         `gorm:"type:varchar(50);uniqueIndex:idx_users_username;not null" json:"username"`
	Email         string         `gorm:"type:varchar(100);uniqueIndex:idx_users_email;not null" json:"email"`
	PasswordHash  string         `gorm:"type:varchar(255);not null" json:"-"`
	Role          UserRole       `gorm:"type:varchar(20);default:user" json:"role"`
	Status        UserStatus     `gorm:"type:smallint;default:1" json:"status"`
	EmailVerified bool           `gorm:"default:false" json:"email_verified"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index:idx_users_deleted_at" json:"-"`
}

func (User) TableName() string {
	return "users"
}
