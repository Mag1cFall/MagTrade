// 密碼雜湊工具
//
// 使用 bcrypt 演算法進行密碼雜湊和驗證
// bcrypt 自帶鹽值，可抵抗彩虹表攻擊
package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword 對密碼進行雜湊
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword 驗證密碼是否匹配
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
