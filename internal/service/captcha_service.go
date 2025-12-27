// 驗證碼服務
//
// 本檔案提供圖形驗證碼和登入安全功能
// 包含：驗證碼生成/驗證、登入失敗記錄、帳號鎖定
// 使用 Redis 儲存驗證碼和失敗計數
package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// CaptchaService 驗證碼服務
type CaptchaService struct {
	redis *redis.Client
}

func NewCaptchaService(rdb *redis.Client) *CaptchaService {
	return &CaptchaService{redis: rdb}
}

// GenerateCaptcha 生成驗證碼
// 返回驗證碼 ID 和驗證碼內容
func (s *CaptchaService) GenerateCaptcha(ctx context.Context, identifier string) (string, string, error) {
	code := s.generateCode(6)
	captchaID := s.generateID()

	key := fmt.Sprintf("captcha:%s", captchaID)
	if err := s.redis.Set(ctx, key, code, 5*time.Minute).Err(); err != nil { // 5 分鐘過期
		return "", "", err
	}

	return captchaID, code, nil
}

// VerifyCaptcha 驗證驗證碼（一次性，驗證後刪除）
func (s *CaptchaService) VerifyCaptcha(ctx context.Context, captchaID, code string) bool {
	key := fmt.Sprintf("captcha:%s", captchaID)
	stored, err := s.redis.Get(ctx, key).Result()
	if err != nil {
		return false
	}

	_ = s.redis.Del(ctx, key) // 無論成功失敗都刪除
	return stored == code
}

// RecordLoginFailure 記錄登入失敗次數
func (s *CaptchaService) RecordLoginFailure(ctx context.Context, identifier string) (int64, error) {
	key := fmt.Sprintf("login_failure:%s", identifier)
	count, err := s.redis.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	if count == 1 { // 首次失敗，設定 15 分鐘過期
		_ = s.redis.Expire(ctx, key, 15*time.Minute)
	}

	return count, nil
}

// GetLoginFailureCount 取得登入失敗次數
func (s *CaptchaService) GetLoginFailureCount(ctx context.Context, identifier string) int64 {
	key := fmt.Sprintf("login_failure:%s", identifier)
	count, _ := s.redis.Get(ctx, key).Int64()
	return count
}

// ClearLoginFailure 清除登入失敗記錄（登入成功時呼叫）
func (s *CaptchaService) ClearLoginFailure(ctx context.Context, identifier string) error {
	key := fmt.Sprintf("login_failure:%s", identifier)
	return s.redis.Del(ctx, key).Err()
}

// IsAccountLocked 檢查帳號是否被鎖定
func (s *CaptchaService) IsAccountLocked(ctx context.Context, identifier string) bool {
	key := fmt.Sprintf("account_locked:%s", identifier)
	exists, _ := s.redis.Exists(ctx, key).Result()
	return exists > 0
}

// LockAccount 鎖定帳號（15 分鐘）
func (s *CaptchaService) LockAccount(ctx context.Context, identifier string) error {
	key := fmt.Sprintf("account_locked:%s", identifier)
	return s.redis.Set(ctx, key, "1", 15*time.Minute).Err()
}

// NeedsCaptcha 檢查是否需要驗證碼（3 次失敗後需要）
func (s *CaptchaService) NeedsCaptcha(ctx context.Context, identifier string) bool {
	return s.GetLoginFailureCount(ctx, identifier) >= 3
}

// generateCode 生成隨機數字驗證碼
func (s *CaptchaService) generateCode(length int) string {
	const digits = "0123456789"
	b := make([]byte, length)
	_, _ = rand.Read(b)
	for i := range b {
		b[i] = digits[int(b[i])%len(digits)]
	}
	return string(b)
}

// generateID 生成驗證碼 ID
func (s *CaptchaService) generateID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
