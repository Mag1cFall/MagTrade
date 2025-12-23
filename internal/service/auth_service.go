package service

import (
	"context"
	"errors"

	"github.com/Mag1cFall/magtrade/internal/config"
	"github.com/Mag1cFall/magtrade/internal/model"
	"github.com/Mag1cFall/magtrade/internal/pkg/utils"
	"github.com/Mag1cFall/magtrade/internal/repository"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserDisabled       = errors.New("user is disabled")
	ErrValidation         = errors.New("validation failed")
	ErrAccountLocked      = errors.New("account is locked")
	ErrCaptchaRequired    = errors.New("captcha required")
	ErrInvalidCaptcha     = errors.New("invalid captcha")
)

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

type RegisterRequest struct {
	Username  string `json:"username" binding:"required,min=3,max=50"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6,max=50"`
	EmailCode string `json:"email_code" binding:"required,len=6"`
}

type LoginRequest struct {
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required"`
	CaptchaID   string `json:"captcha_id"`
	CaptchaCode string `json:"captcha_code"`
}

type SendCodeRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type TokenResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresIn    int       `json:"expires_in"`
	User         *UserInfo `json:"user"`
}

type UserInfo struct {
	ID       int64            `json:"id"`
	Username string           `json:"username"`
	Email    string           `json:"email"`
	Role     model.UserRole   `json:"role"`
	Status   model.UserStatus `json:"status"`
}

func (s *AuthService) SendEmailCode(ctx context.Context, email string) error {
	exists, err := s.userRepo.ExistsByEmail(ctx, email)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("email already registered")
	}

	return s.emailService.SendEmailCode(ctx, email)
}

func (s *AuthService) Register(ctx context.Context, req *RegisterRequest) (*TokenResponse, error) {
	if !s.emailService.VerifyEmailCode(ctx, req.Email, req.EmailCode) {
		return nil, ErrInvalidEmailCode
	}

	exists, err := s.userRepo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("username already exists")
	}

	exists, err = s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email already exists")
	}

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

func (s *AuthService) Login(ctx context.Context, req *LoginRequest) (*TokenResponse, error) {
	identifier := req.Username

	if s.captchaService.IsAccountLocked(ctx, identifier) {
		return nil, ErrAccountLocked
	}

	needsCaptcha := s.captchaService.NeedsCaptcha(ctx, identifier)
	if needsCaptcha {
		if req.CaptchaID == "" || req.CaptchaCode == "" {
			return nil, ErrCaptchaRequired
		}
		if !s.captchaService.VerifyCaptcha(ctx, req.CaptchaID, req.CaptchaCode) {
			return nil, ErrInvalidCaptcha
		}
	}

	user, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			s.handleLoginFailure(ctx, identifier)
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		s.handleLoginFailure(ctx, identifier)
		return nil, ErrInvalidCredentials
	}

	if user.Status == model.UserStatusDisabled {
		return nil, ErrUserDisabled
	}

	_ = s.captchaService.ClearLoginFailure(ctx, identifier)

	return s.generateTokenResponse(user)
}

func (s *AuthService) handleLoginFailure(ctx context.Context, identifier string) {
	count, _ := s.captchaService.RecordLoginFailure(ctx, identifier)
	if count >= 5 {
		_ = s.captchaService.LockAccount(ctx, identifier)
	}
}

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
