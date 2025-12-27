// 認證業務服務
//
// 本檔案處理使用者認證相關業務邏輯
// 包含：註冊、登入、Token 刷新、登入失敗鎖定
// 整合驗證碼服務和郵件服務
package service

import (
	"context"
	"errors"

	"github.com/Mag1cFall/magtrade/internal/config"
	"github.com/Mag1cFall/magtrade/internal/model"
	"github.com/Mag1cFall/magtrade/internal/pkg/utils"
	"github.com/Mag1cFall/magtrade/internal/repository"
)

// 業務錯誤定義
var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserDisabled       = errors.New("user is disabled")
	ErrValidation         = errors.New("validation failed")
	ErrAccountLocked      = errors.New("account is locked")
	ErrCaptchaRequired    = errors.New("captcha required")
	ErrInvalidCaptcha     = errors.New("invalid captcha")
)

// AuthService 認證服務
type AuthService struct {
	userRepo       *repository.UserRepository
	jwtCfg         *config.JWTConfig
	captchaService *CaptchaService
	emailService   *EmailService
}

func NewAuthService(jwtCfg *config.JWTConfig, captchaSvc *CaptchaService, emailSvc *EmailService) *AuthService {
	return &AuthService{
		userRepo:       repository.NewUserRepository(),
		jwtCfg:         jwtCfg,
		captchaService: captchaSvc,
		emailService:   emailSvc,
	}
}

// RegisterRequest 註冊請求
type RegisterRequest struct {
	Username  string `json:"username" binding:"required,min=3,max=50"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6,max=50"`
	EmailCode string `json:"email_code" binding:"required,len=6"` // 郵件驗證碼
}

// LoginRequest 登入請求
type LoginRequest struct {
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required"`
	CaptchaID   string `json:"captcha_id"`   // 驗證碼 ID（多次失敗後需要）
	CaptchaCode string `json:"captcha_code"` // 驗證碼
}

// SendCodeRequest 發送驗證碼請求
type SendCodeRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// TokenResponse 登入成功回應
type TokenResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresIn    int       `json:"expires_in"` // 過期時間（秒）
	User         *UserInfo `json:"user"`
}

// UserInfo 使用者基本資訊
type UserInfo struct {
	ID       int64            `json:"id"`
	Username string           `json:"username"`
	Email    string           `json:"email"`
	Role     model.UserRole   `json:"role"`
	Status   model.UserStatus `json:"status"`
}

// SendEmailCode 發送郵件驗證碼
func (s *AuthService) SendEmailCode(ctx context.Context, email string) error {
	// 檢查郵件是否已註冊
	exists, err := s.userRepo.ExistsByEmail(ctx, email)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("email already registered")
	}

	return s.emailService.SendEmailCode(ctx, email)
}

// Register 使用者註冊
func (s *AuthService) Register(ctx context.Context, req *RegisterRequest) (*TokenResponse, error) {
	// 驗證郵件驗證碼
	if !s.emailService.VerifyEmailCode(ctx, req.Email, req.EmailCode) {
		return nil, ErrInvalidEmailCode
	}

	// 檢查使用者名稱是否存在
	exists, err := s.userRepo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("username already exists")
	}

	// 檢查郵件是否存在
	exists, err = s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email already exists")
	}

	// 密碼雜湊
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username:      req.Username,
		Email:         req.Email,
		PasswordHash:  passwordHash,
		Role:          model.UserRoleUser,
		Status:        model.UserStatusActive,
		EmailVerified: true,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return s.generateTokenResponse(user)
}

// Login 使用者登入
func (s *AuthService) Login(ctx context.Context, req *LoginRequest) (*TokenResponse, error) {
	identifier := req.Username

	// 檢查帳號是否被鎖定
	if s.captchaService.IsAccountLocked(ctx, identifier) {
		return nil, ErrAccountLocked
	}

	// 檢查是否需要驗證碼（3 次失敗後）
	needsCaptcha := s.captchaService.NeedsCaptcha(ctx, identifier)
	if needsCaptcha {
		if req.CaptchaID == "" || req.CaptchaCode == "" {
			return nil, ErrCaptchaRequired
		}
		if !s.captchaService.VerifyCaptcha(ctx, req.CaptchaID, req.CaptchaCode) {
			return nil, ErrInvalidCaptcha
		}
	}

	// 查詢使用者
	user, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			s.handleLoginFailure(ctx, identifier)
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	// 驗證密碼
	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		s.handleLoginFailure(ctx, identifier)
		return nil, ErrInvalidCredentials
	}

	// 檢查帳號狀態
	if user.Status == model.UserStatusDisabled {
		return nil, ErrUserDisabled
	}

	// 登入成功，清除失敗記錄
	_ = s.captchaService.ClearLoginFailure(ctx, identifier)

	return s.generateTokenResponse(user)
}

// handleLoginFailure 處理登入失敗（記錄失敗次數，達到閾值則鎖定）
func (s *AuthService) handleLoginFailure(ctx context.Context, identifier string) {
	count, _ := s.captchaService.RecordLoginFailure(ctx, identifier)
	if count >= 5 { // 5 次失敗鎖定帳號
		_ = s.captchaService.LockAccount(ctx, identifier)
	}
}

// RefreshToken 刷新 Token
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	claims, err := utils.ValidateRefreshToken(refreshToken, s.jwtCfg.Secret)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	if user.Status == model.UserStatusDisabled {
		return nil, ErrUserDisabled
	}

	return s.generateTokenResponse(user)
}

// GetUserByID 根據 ID 取得使用者資訊
func (s *AuthService) GetUserByID(ctx context.Context, userID int64) (*UserInfo, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &UserInfo{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
		Status:   user.Status,
	}, nil
}

// generateTokenResponse 生成 Token 回應
func (s *AuthService) generateTokenResponse(user *model.User) (*TokenResponse, error) {
	accessToken, refreshToken, err := utils.GenerateTokenPair(
		user.ID,
		user.Username,
		string(user.Role),
		s.jwtCfg,
	)
	if err != nil {
		return nil, err
	}

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int(s.jwtCfg.AccessTokenExpire.Seconds()),
		User: &UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role,
			Status:   user.Status,
		},
	}, nil
}
