// 密碼雜湊工具單元測試
//
// 測試覆蓋：
// - HashPassword: bcrypt 雜湊生成（非空、非明文、鹽值唯一性）
// - CheckPassword: 密碼驗證（正確密碼、錯誤密碼、空密碼、空雜湊）
package utils

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "testPassword123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}

	if hash == "" {
		t.Error("HashPassword() returned empty hash")
	}

	if hash == password {
		t.Error("HashPassword() should not return plain password")
	}

	hash2, _ := HashPassword(password)
	if hash == hash2 {
		t.Error("HashPassword() should generate different hashes for same password (salt)")
	}
}

func TestCheckPassword(t *testing.T) {
	password := "testPassword123"
	hash, _ := HashPassword(password)

	tests := []struct {
		name     string
		password string
		hash     string
		want     bool
	}{
		{
			name:     "correct password",
			password: password,
			hash:     hash,
			want:     true,
		},
		{
			name:     "wrong password",
			password: "wrongPassword",
			hash:     hash,
			want:     false,
		},
		{
			name:     "empty password",
			password: "",
			hash:     hash,
			want:     false,
		},
		{
			name:     "empty hash",
			password: password,
			hash:     "",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckPassword(tt.password, tt.hash); got != tt.want {
				t.Errorf("CheckPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}
