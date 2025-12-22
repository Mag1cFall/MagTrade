package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type CaptchaService struct {
	redis *redis.Client
}

func NewCaptchaService(rdb *redis.Client) *CaptchaService {
	return &CaptchaService{redis: rdb}
}

func (s *CaptchaService) GenerateCaptcha(ctx context.Context, identifier string) (string, string, error) {
	code := s.generateCode(6)
	captchaID := s.generateID()

	key := fmt.Sprintf("captcha:%s", captchaID)
	if err := s.redis.Set(ctx, key, code, 5*time.Minute).Err(); err != nil {
		return "", "", err
	}

	return captchaID, code, nil
}

func (s *CaptchaService) VerifyCaptcha(ctx context.Context, captchaID, code string) bool {
	key := fmt.Sprintf("captcha:%s", captchaID)
	stored, err := s.redis.Get(ctx, key).Result()
	if err != nil {
		return false
	}

	_ = s.redis.Del(ctx, key)
	return stored == code
}

func (s *CaptchaService) RecordLoginFailure(ctx context.Context, identifier string) (int64, error) {
	key := fmt.Sprintf("login_failure:%s", identifier)
	count, err := s.redis.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	if count == 1 {
		_ = s.redis.Expire(ctx, key, 15*time.Minute)
	}

	return count, nil
}

func (s *CaptchaService) GetLoginFailureCount(ctx context.Context, identifier string) int64 {
	key := fmt.Sprintf("login_failure:%s", identifier)
	count, _ := s.redis.Get(ctx, key).Int64()
	return count
}

func (s *CaptchaService) ClearLoginFailure(ctx context.Context, identifier string) error {
	key := fmt.Sprintf("login_failure:%s", identifier)
	return s.redis.Del(ctx, key).Err()
}

func (s *CaptchaService) IsAccountLocked(ctx context.Context, identifier string) bool {
	key := fmt.Sprintf("account_locked:%s", identifier)
	exists, _ := s.redis.Exists(ctx, key).Result()
	return exists > 0
}

func (s *CaptchaService) LockAccount(ctx context.Context, identifier string) error {
	key := fmt.Sprintf("account_locked:%s", identifier)
	return s.redis.Set(ctx, key, "1", 15*time.Minute).Err()
}

func (s *CaptchaService) NeedsCaptcha(ctx context.Context, identifier string) bool {
	return s.GetLoginFailureCount(ctx, identifier) >= 3
}

func (s *CaptchaService) generateCode(length int) string {
	const digits = "0123456789"
	b := make([]byte, length)
	_, _ = rand.Read(b)
	for i := range b {
		b[i] = digits[int(b[i])%len(digits)]
	}
	return string(b)
}

func (s *CaptchaService) generateID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
