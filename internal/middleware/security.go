// 安全防護中間件
//
// 本檔案提供多種安全防護功能：
// Security: HTTP 安全標頭
// XSSFilter: XSS 攻擊過濾
// SQLInjectionFilter: SQL 注入過濾
// RequestSizeLimit: 請求大小限制
package middleware

import (
	"strings"

	"github.com/Mag1cFall/magtrade/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

// Security HTTP 安全標頭中間件
func Security() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")                  // 禁止 MIME 類型嗅探
		c.Header("X-Frame-Options", "DENY")                            // 禁止 iframe 嵌入（防點擊劫持）
		c.Header("X-XSS-Protection", "1; mode=block")                  // 啟用瀏覽器 XSS 過濾
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin") // 限制 Referrer 洩漏

		// Swagger 頁面需要載入外部資源，放寬 CSP
		if c.Request.URL.Path == "/swagger" || c.Request.URL.Path == "/swagger.yaml" {
			c.Header("Content-Security-Policy", "default-src 'self' 'unsafe-inline' 'unsafe-eval' https://unpkg.com; img-src 'self' data: https:; style-src 'self' 'unsafe-inline' https://unpkg.com;")
		} else {
			c.Header("Content-Security-Policy", "default-src 'self'") // 只允許同源資源
		}

		c.Next()
	}
}

// XSSFilter XSS 攻擊過濾中間件
// 檢測 URL Query 中的危險標籤
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

// SQLInjectionFilter SQL 注入過濾中間件
// 檢測 URL Query 中的危險模式
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

// RequestSizeLimit 請求大小限制中間件
// 防止超大請求消耗伺服器資源
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
