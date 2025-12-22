package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/smtp"
	"time"

	"github.com/Mag1cFall/magtrade/internal/config"
	"github.com/Mag1cFall/magtrade/internal/repository"
	"github.com/redis/go-redis/v9"
)

var (
	ErrEmailNotVerified = errors.New("email not verified")
	ErrInvalidToken     = errors.New("invalid or expired token")
)

type EmailConfig struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	FromAddress  string
	FromName     string
}

type EmailService struct {
	redis    *redis.Client
	userRepo *repository.UserRepository
	cfg      *EmailConfig
}

func NewEmailService(rdb *redis.Client, cfg *EmailConfig) *EmailService {
	return &EmailService{
		redis:    rdb,
		userRepo: repository.NewUserRepository(),
		cfg:      cfg,
	}
}

func (s *EmailService) SendVerificationEmail(ctx context.Context, userID int64, email string) error {
	token := s.generateToken()
	key := fmt.Sprintf("email_verify:%s", token)

	data := fmt.Sprintf("%d:%s", userID, email)
	if err := s.redis.Set(ctx, key, data, 24*time.Hour).Err(); err != nil {
		return err
	}

	if s.cfg == nil || s.cfg.SMTPHost == "" {
		return nil
	}

	subject := "验证您的邮箱 - MagTrade"
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head><meta charset="utf-8"></head>
<body>
<h2>欢迎注册 MagTrade</h2>
<p>请点击下方链接验证您的邮箱：</p>
<p><a href="http://localhost:8080/api/v1/auth/verify-email?token=%s">验证邮箱</a></p>
<p>链接24小时内有效。</p>
</body>
</html>
`, token)

	return s.sendEmail(email, subject, body)
}

func (s *EmailService) VerifyEmail(ctx context.Context, token string) error {
	key := fmt.Sprintf("email_verify:%s", token)
	data, err := s.redis.Get(ctx, key).Result()
	if err != nil {
		return ErrInvalidToken
	}

	var userID int64
	var email string
	_, err = fmt.Sscanf(data, "%d:%s", &userID, &email)
	if err != nil {
		return ErrInvalidToken
	}

	if err := s.userRepo.UpdateEmailVerified(ctx, userID, true); err != nil {
		return err
	}

	_ = s.redis.Del(ctx, key)
	return nil
}

func (s *EmailService) sendEmail(to, subject, body string) error {
	if s.cfg == nil {
		return nil
	}

	from := s.cfg.FromAddress
	msg := fmt.Sprintf("From: %s <%s>\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=utf-8\r\n"+
		"\r\n"+
		"%s", s.cfg.FromName, from, to, subject, body)

	auth := smtp.PlainAuth("", s.cfg.SMTPUser, s.cfg.SMTPPassword, s.cfg.SMTPHost)
	addr := fmt.Sprintf("%s:%d", s.cfg.SMTPHost, s.cfg.SMTPPort)

	return smtp.SendMail(addr, auth, from, []string{to}, []byte(msg))
}

func (s *EmailService) generateToken() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func LoadEmailConfig(cfg *config.Config) *EmailConfig {
	return nil
}
