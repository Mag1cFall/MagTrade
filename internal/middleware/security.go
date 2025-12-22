package middleware

import (
	"strings"

	"github.com/Mag1cFall/magtrade/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

func Security() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		if c.Request.URL.Path == "/swagger" || c.Request.URL.Path == "/swagger.yaml" {
			c.Header("Content-Security-Policy", "default-src 'self' 'unsafe-inline' 'unsafe-eval' https://unpkg.com; img-src 'self' data: https:; style-src 'self' 'unsafe-inline' https://unpkg.com;")
		} else {
			c.Header("Content-Security-Policy", "default-src 'self'")
		}

		c.Next()
	}
}

func XSSFilter() gin.HandlerFunc {
	dangerousTags := []string{"<script", "javascript:", "onerror=", "onload=", "onclick=", "onmouseover="}

	return func(c *gin.Context) {
		query := c.Request.URL.RawQuery
		queryLower := strings.ToLower(query)

		for _, tag := range dangerousTags {
			if strings.Contains(queryLower, tag) {
				response.BadRequest(c, "potential XSS detected")
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

func SQLInjectionFilter() gin.HandlerFunc {
	dangerousPatterns := []string{
		"' or '1'='1",
		"'; drop table",
		"union select",
		"--",
		";--",
		"/*",
		"*/",
		"@@version",
		"char(",
		"nchar(",
		"varchar(",
		"nvarchar(",
		"cast(",
		"convert(",
		"exec(",
		"execute(",
	}

	return func(c *gin.Context) {
		query := strings.ToLower(c.Request.URL.RawQuery)

		for _, pattern := range dangerousPatterns {
			if strings.Contains(query, pattern) {
				response.BadRequest(c, "potential SQL injection detected")
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

func RequestSizeLimit(maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.ContentLength > maxSize {
			response.BadRequest(c, "request body too large")
			c.Abort()
			return
		}

		c.Next()
	}
}
