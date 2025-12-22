package handler

import (
	"github.com/Mag1cFall/magtrade/internal/config"
	"github.com/Mag1cFall/magtrade/internal/pkg/response"
	"github.com/Mag1cFall/magtrade/internal/pkg/utils"
	"github.com/Mag1cFall/magtrade/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
	jwtCfg      *config.JWTConfig
}

func NewAuthHandler(jwtCfg *config.JWTConfig) *AuthHandler {
	return &AuthHandler{
		authService: service.NewAuthService(jwtCfg),
		jwtCfg:      jwtCfg,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
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
		if err == service.ErrInvalidCredentials {
			response.Unauthorized(c, "用户名或密码错误")
			return
		}
		if err == service.ErrUserDisabled {
			response.Forbidden(c, "账号已被禁用")
			return
		}
		response.InternalError(c, err.Error())
		return
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
