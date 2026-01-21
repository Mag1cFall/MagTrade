// 使用者模型單元測試
//
// 測試覆蓋：
// - User.TableName: GORM 資料表名稱
// - UserStatus 常量: 啟用/停用狀態值驗證
// - UserRole 常量: user/admin 角色值驗證
// - User 欄位: 模型欄位賦值與讀取
package model

import (
	"testing"
)

func TestUser_TableName(t *testing.T) {
	u := User{}
	if got := u.TableName(); got != "users" {
		t.Errorf("User.TableName() = %v, want %v", got, "users")
	}
}

func TestUserStatus_Constants(t *testing.T) {
	tests := []struct {
		name   string
		status UserStatus
		want   int8
	}{
		{"disabled", UserStatusDisabled, 0},
		{"active", UserStatusActive, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if int8(tt.status) != tt.want {
				t.Errorf("UserStatus %s = %v, want %v", tt.name, int8(tt.status), tt.want)
			}
		})
	}
}

func TestUserRole_Constants(t *testing.T) {
	tests := []struct {
		name string
		role UserRole
		want string
	}{
		{"user role", UserRoleUser, "user"},
		{"admin role", UserRoleAdmin, "admin"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.role) != tt.want {
				t.Errorf("UserRole %s = %v, want %v", tt.name, string(tt.role), tt.want)
			}
		})
	}
}

func TestUser_Fields(t *testing.T) {
	u := User{
		ID:            1,
		Username:      "testuser",
		Email:         "test@example.com",
		PasswordHash:  "hashedpassword",
		Role:          UserRoleUser,
		Status:        UserStatusActive,
		EmailVerified: true,
	}

	if u.ID != 1 {
		t.Errorf("User.ID = %v, want %v", u.ID, 1)
	}
	if u.Username != "testuser" {
		t.Errorf("User.Username = %v, want %v", u.Username, "testuser")
	}
	if u.Email != "test@example.com" {
		t.Errorf("User.Email = %v, want %v", u.Email, "test@example.com")
	}
	if u.PasswordHash != "hashedpassword" {
		t.Errorf("User.PasswordHash = %v, want %v", u.PasswordHash, "hashedpassword")
	}
	if u.Role != UserRoleUser {
		t.Errorf("User.Role = %v, want %v", u.Role, UserRoleUser)
	}
	if u.Status != UserStatusActive {
		t.Errorf("User.Status = %v, want %v", u.Status, UserStatusActive)
	}
	if !u.EmailVerified {
		t.Errorf("User.EmailVerified = %v, want %v", u.EmailVerified, true)
	}
}

func TestUser_AdminRole(t *testing.T) {
	u := User{
		Role: UserRoleAdmin,
	}

	if u.Role != UserRoleAdmin {
		t.Errorf("User.Role = %v, want %v", u.Role, UserRoleAdmin)
	}
	if string(u.Role) != "admin" {
		t.Errorf("User.Role string = %v, want %v", string(u.Role), "admin")
	}
}
