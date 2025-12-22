package middleware

import (
	"strings"

	"github.com/Mag1cFall/magtrade/internal/config"
	"github.com/Mag1cFall/magtrade/internal/model"
	"github.com/Mag1cFall/magtrade/internal/pkg/response"
	"github.com/Mag1cFall/magtrade/internal/pkg/utils"
	"github.com/gin-gonic/gin"
)

const (
	AuthorizationHeader = "Authorization"
	BearerPrefix        = "Bearer "
	ContextUserID       = "user_id"
	ContextUsername     = "username"
	ContextUserRole     = "user_role"
)

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

		c.Set(ContextUserID, claims.UserID)
		c.Set(ContextUsername, claims.Username)
		c.Set(ContextUserRole, claims.Role)

		c.Next()
	}
}

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

func GetUserID(c *gin.Context) int64 {
	userID, exists := c.Get(ContextUserID)
	if !exists {
		return 0
	}
	return userID.(int64)
}

func GetUsername(c *gin.Context) string {
	username, exists := c.Get(ContextUsername)
	if !exists {
		return ""
	}
	return username.(string)
}

func GetUserRole(c *gin.Context) string {
	role, exists := c.Get(ContextUserRole)
	if !exists {
		return ""
	}
	return role.(string)
}
