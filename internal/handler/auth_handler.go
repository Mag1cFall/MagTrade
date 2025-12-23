package handler

import (
	"github.com/Mag1cFall/magtrade/internal/cache"
	"github.com/Mag1cFall/magtrade/internal/config"
	"github.com/Mag1cFall/magtrade/internal/pkg/response"
	"github.com/Mag1cFall/magtrade/internal/pkg/utils"
	"github.com/Mag1cFall/magtrade/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService    *service.AuthService
	captchaService *service.CaptchaService
	jwtCfg         *config.JWTConfig
}

func NewAuthHandler(jwtCfg *config.JWTConfig, emailCfg *config.EmailConfig) *AuthHandler {
	rdb := cache.GetClient()
	captchaSvc := service.NewCaptchaService(rdb)
	emailSvc := service.NewEmailService(rdb, emailCfg)

	return &AuthHandler{
		authService:    service.NewAuthService(jwtCfg, captchaSvc, emailSvc),
		captchaService: captchaSvc,
		jwtCfg:         jwtCfg,
	}
}

func (h *AuthHandler) SendEmailCode(c *gin.Context) {
	var req service.SendCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.authService.SendEmailCode(c.Request.Context(), req.Email); err != nil {
		if err == service.ErrEmailCodeTooFrequent {
			response.BadRequest(c, "请60秒后再试")
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "验证码已发送"})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		if err == service.ErrInvalidEmailCode {
			response.BadRequest(c, "验证码错误或已过期")
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, result)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		switch err {
		case service.ErrInvalidCredentials:
			identifier := req.Username
			needsCaptcha := h.captchaService.NeedsCaptcha(c.Request.Context(), identifier)
			response.Unauthorized(c, "用户名或密码错误", gin.H{"needs_captcha": needsCaptcha})
			return
		case service.ErrUserDisabled:
			response.Forbidden(c, "账号已被禁用")
			return
		case service.ErrAccountLocked:
			response.Forbidden(c, "账号已被锁定，请15分钟后再试")
			return
		case service.ErrCaptchaRequired:
			response.BadRequest(c, "请输入验证码")
			return
		case service.ErrInvalidCaptcha:
			response.BadRequest(c, "验证码错误")
			return
		default:
			response.InternalError(c, err.Error())
			return
		}
	}

	response.Success(c, result)
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		if err == utils.ErrExpiredToken {
			response.Unauthorized(c, "refresh token已过期，请重新登录")
			return
		}
		response.Unauthorized(c, "invalid refresh token")
		return
	}

	response.Success(c, result)
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "authentication required")
		return
	}

	user, err := h.authService.GetUserByID(c.Request.Context(), userID.(int64))
	if err != nil {
		response.NotFound(c, "user not found")
		return
	}

	response.Success(c, user)
}
