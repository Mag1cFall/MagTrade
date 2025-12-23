package service

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"errors"
	"fmt"
	"net/smtp"
	"time"

	"github.com/Mag1cFall/magtrade/internal/config"
	"github.com/redis/go-redis/v9"
)

var (
	ErrEmailCodeTooFrequent = errors.New("please wait before requesting a new code")
	ErrInvalidEmailCode     = errors.New("invalid or expired verification code")
)

type EmailService struct {
	redis *redis.Client
	cfg   *config.EmailConfig
}

func NewEmailService(rdb *redis.Client, cfg *config.EmailConfig) *EmailService {
	return &EmailService{
		redis: rdb,
		cfg:   cfg,
	}
}

func (s *EmailService) SendEmailCode(ctx context.Context, email string) error {
	cooldownKey := fmt.Sprintf("email_code_cooldown:%s", email)
	exists, _ := s.redis.Exists(ctx, cooldownKey).Result()
	if exists > 0 {
		return ErrEmailCodeTooFrequent
	}

	code := s.generateCode(6)
	codeKey := fmt.Sprintf("email_code:%s", email)

	if err := s.redis.Set(ctx, codeKey, code, 15*time.Minute).Err(); err != nil {
		return err
	}

	if err := s.redis.Set(ctx, cooldownKey, "1", 60*time.Second).Err(); err != nil {
		return err
	}

	if s.cfg == nil || s.cfg.SMTPHost == "" {
		fmt.Printf("[DEV MODE] Email verification code for %s: %s\n", email, code)
		return nil
	}

	subject := "您的验证码 - MagTrade"
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head><meta charset="utf-8"></head>
<body style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;">
<div style="background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); padding: 30px; border-radius: 10px; text-align: center;">
<h1 style="color: white; margin: 0;">MagTrade</h1>
<p style="color: rgba(255,255,255,0.9); margin-top: 10px;">高并发分布式秒杀系统</p>
</div>
<div style="padding: 30px 0;">
<p>您好，</p>
<p>您正在注册 MagTrade 账号，验证码为：</p>
<div style="background: #f5f5f5; padding: 20px; text-align: center; border-radius: 8px; margin: 20px 0;">
<span style="font-size: 32px; font-weight: bold; letter-spacing: 8px; color: #667eea;">%s</span>
</div>
<p>验证码 <strong>15分钟</strong> 内有效，请勿泄露给他人。</p>
<p style="color: #999; font-size: 12px; margin-top: 30px;">如果这不是您的操作，请忽略此邮件。</p>
</div>
</body>
</html>
`, code)

	return s.sendEmail(email, subject, body)
}

func (s *EmailService) VerifyEmailCode(ctx context.Context, email, code string) bool {
	codeKey := fmt.Sprintf("email_code:%s", email)
	stored, err := s.redis.Get(ctx, codeKey).Result()
	if err != nil {
		return false
	}

	if stored == code {
		_ = s.redis.Del(ctx, codeKey)
		return true
	}
	return false
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

	addr := fmt.Sprintf("%s:%d", s.cfg.SMTPHost, s.cfg.SMTPPort)

	tlsConfig := &tls.Config{
		ServerName: s.cfg.SMTPHost,
	}

	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, s.cfg.SMTPHost)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Close()

	auth := smtp.PlainAuth("", s.cfg.SMTPUser, s.cfg.SMTPPassword, s.cfg.SMTPHost)
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("SMTP auth failed: %w", err)
	}

	if err := client.Mail(from); err != nil {
		return fmt.Errorf("SMTP MAIL command failed: %w", err)
	}

	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("SMTP RCPT command failed: %w", err)
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("SMTP DATA command failed: %w", err)
	}

	_, err = w.Write([]byte(msg))
	if err != nil {
		return fmt.Errorf("failed to write email body: %w", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("failed to close email writer: %w", err)
	}

	return client.Quit()
}

func (s *EmailService) generateCode(length int) string {
	const digits = "0123456789"
	b := make([]byte, length)
	_, _ = rand.Read(b)
	for i := range b {
		b[i] = digits[int(b[i])%len(digits)]
	}
	return string(b)
}
