// JWT 認證中間件
//
// 本檔案提供 JWT Token 驗證和使用者資訊提取功能
// 支援 Bearer Token 格式，驗證後將使用者資訊存入 Context
package middleware

import (
	"strings"

	"github.com/Mag1cFall/magtrade/internal/config"
	"github.com/Mag1cFall/magtrade/internal/model"
	"github.com/Mag1cFall/magtrade/internal/pkg/response"
	"github.com/Mag1cFall/magtrade/internal/pkg/utils"
	"github.com/gin-gonic/gin"
)

// Context Key 常量定義
const (
	AuthorizationHeader = "Authorization"
	BearerPrefix        = "Bearer "
	ContextUserID       = "user_id"
	ContextUsername     = "username"
	ContextUserRole     = "user_role"
)

// Auth JWT 認證中間件
// 從 Authorization Header 取得並驗證 Token
// 驗證成功後將 user_id/username/role 存入 Context
func Auth(jwtCfg *config.JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(AuthorizationHeader)
		if authHeader == "" {
			response.Unauthorized(c, "missing authorization header")
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, BearerPrefix) {
			response.Unauthorized(c, "invalid authorization format")
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, BearerPrefix)

		claims, err := utils.ValidateAccessToken(tokenString, jwtCfg.Secret)
		if err != nil {
			if err == utils.ErrExpiredToken {
				response.Unauthorized(c, "token has expired")
			} else {
				response.Unauthorized(c, "invalid token")
			}
			c.Abort()
			return
		}

		// 將使用者資訊存入 Context，供後續 Handler 使用
		c.Set(ContextUserID, claims.UserID)
		c.Set(ContextUsername, claims.Username)
		c.Set(ContextUserRole, claims.Role)

		c.Next()
	}
}

// AdminOnly 管理員權限檢查中間件
// 必須搭配 Auth 中間件使用
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get(ContextUserRole)
		if !exists {
			response.Unauthorized(c, "authentication required")
			c.Abort()
			return
		}

		if role != string(model.UserRoleAdmin) {
			response.Forbidden(c, "admin access required")
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetUserID 從 Context 取得當前使用者 ID
func GetUserID(c *gin.Context) int64 {
	userID, exists := c.Get(ContextUserID)
	if !exists {
		return 0
	}
	return userID.(int64)
}

// GetUsername 從 Context 取得當前使用者名稱
func GetUsername(c *gin.Context) string {
	username, exists := c.Get(ContextUsername)
	if !exists {
		return ""
	}
	return username.(string)
}

// GetUserRole 從 Context 取得當前使用者角色
func GetUserRole(c *gin.Context) string {
	role, exists := c.Get(ContextUserRole)
	if !exists {
		return ""
	}
	return role.(string)
}
