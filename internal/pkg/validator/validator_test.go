// 資料驗證工具單元測試
//
// 測試覆蓋：
// - ValidateUsername: 使用者名稱格式驗證（長度、字元規則）
// - ValidateEmail: 郵箱格式驗證
// - ValidatePassword: 密碼強度驗證
// - ValidatePhone: 手機號碼格式驗證（中國大陸）
// - ValidateQuantity: 購買數量範圍驗證
// - ValidateFlashSaleID: 秒殺活動 ID 驗證
// - ValidateOrderNo: 訂單號格式驗證
// - ValidateSessionID: 會話 ID 驗證
// - ValidateMessage: 消息內容長度驗證
// - ValidationError: 錯誤結構體方法測試
package validator

import (
	"testing"
)

func TestValidateUsername(t *testing.T) {
	tests := []struct {
		name     string
		username string
		wantErr  bool
	}{
		{"valid username", "testuser", false},
		{"valid with underscore", "test_user123", false},
		{"too short", "ab", true},
		{"too long", "abcdefghijklmnopqrstuvwxyz12345", true},
		{"starts with number", "1testuser", true},
		{"contains special char", "test@user", true},
		{"valid min length", "abc", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUsername(tt.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateUsername(%q) error = %v, wantErr %v", tt.username, err, tt.wantErr)
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{"valid email", "test@example.com", false},
		{"valid with plus", "test+tag@example.com", false},
		{"empty email", "", true},
		{"no at symbol", "testexample.com", true},
		{"no domain", "test@", true},
		{"invalid format", "test@@example.com", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmail(tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateEmail(%q) error = %v, wantErr %v", tt.email, err, tt.wantErr)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{"valid password", "password123", false},
		{"valid complex", "Password123!", false},
		{"too short", "12345", true},
		{"exactly 6 chars", "123456", false},
		{"empty password", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePassword(%q) error = %v, wantErr %v", tt.password, err, tt.wantErr)
			}
		})
	}
}

func TestValidateQuantity(t *testing.T) {
	tests := []struct {
		name     string
		quantity int
		wantErr  bool
	}{
		{"valid quantity", 1, false},
		{"max valid", 10, false},
		{"zero", 0, true},
		{"negative", -1, true},
		{"too large", 11, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateQuantity(tt.quantity)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateQuantity(%d) error = %v, wantErr %v", tt.quantity, err, tt.wantErr)
			}
		})
	}
}

func TestValidatePhone(t *testing.T) {
	tests := []struct {
		name    string
		phone   string
		wantErr bool
	}{
		{"valid phone", "13812345678", false},
		{"valid phone 2", "15912345678", false},
		{"valid phone 3", "18888888888", false},
		{"empty phone (optional)", "", false},
		{"too short", "1381234567", true},
		{"too long", "138123456789", true},
		{"invalid prefix", "12812345678", true},
		{"contains letter", "1381234567a", true},
		{"landline number", "02112345678", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePhone(tt.phone)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePhone(%q) error = %v, wantErr %v", tt.phone, err, tt.wantErr)
			}
		})
	}
}

func TestValidateFlashSaleID(t *testing.T) {
	tests := []struct {
		name    string
		id      int64
		wantErr bool
	}{
		{"valid id", 1, false},
		{"valid large id", 9999999, false},
		{"zero id", 0, true},
		{"negative id", -1, true},
		{"negative large", -9999, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFlashSaleID(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateFlashSaleID(%d) error = %v, wantErr %v", tt.id, err, tt.wantErr)
			}
		})
	}
}

func TestValidateOrderNo(t *testing.T) {
	tests := []struct {
		name    string
		orderNo string
		wantErr bool
	}{
		{"valid order no", "FS1234567890", false},
		{"valid 10 chars", "1234567890", false},
		{"valid 50 chars", "12345678901234567890123456789012345678901234567890", false},
		{"empty order no", "", true},
		{"too short", "123456789", true},
		{"too long 51 chars", "123456789012345678901234567890123456789012345678901", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateOrderNo(tt.orderNo)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateOrderNo(%q) error = %v, wantErr %v", tt.orderNo, err, tt.wantErr)
			}
		})
	}
}

func TestValidateSessionID(t *testing.T) {
	tests := []struct {
		name      string
		sessionID string
		wantErr   bool
	}{
		{"valid session id", "abc123", false},
		{"valid uuid format", "550e8400-e29b-41d4-a716-446655440000", false},
		{"valid 64 chars", "1234567890123456789012345678901234567890123456789012345678901234", false},
		{"empty session id", "", true},
		{"too long 65 chars", "12345678901234567890123456789012345678901234567890123456789012345", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSessionID(tt.sessionID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSessionID(%q) error = %v, wantErr %v", tt.sessionID, err, tt.wantErr)
			}
		})
	}
}

func TestValidateMessage(t *testing.T) {
	tests := []struct {
		name    string
		message string
		wantErr bool
	}{
		{"valid message", "Hello, world!", false},
		{"valid single char", "H", false},
		{"valid 2000 chars", string(make([]byte, 2000)), false},
		{"empty message", "", true},
		{"too long 2001 chars", string(make([]byte, 2001)), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid 2000 chars" {
				tt.message = string(make([]rune, 2000))
				for i := range tt.message {
					tt.message = tt.message[:i] + "a" + tt.message[i+1:]
				}
				tt.message = repeatString("a", 2000)
			}
			if tt.name == "too long 2001 chars" {
				tt.message = repeatString("a", 2001)
			}
			err := ValidateMessage(tt.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func repeatString(s string, count int) string {
	result := make([]byte, len(s)*count)
	for i := 0; i < count; i++ {
		copy(result[i*len(s):], s)
	}
	return string(result)
}

func TestValidationError_Error(t *testing.T) {
	err := &ValidationError{
		Field:   "username",
		Message: "用户名至少3个字符",
	}

	if err.Error() != "用户名至少3个字符" {
		t.Errorf("ValidationError.Error() = %v, want %v", err.Error(), "用户名至少3个字符")
	}
}
