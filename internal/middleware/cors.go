// CORS 跨域資源共享中間件
//
// 允許前端跨域訪問 API
// 生產環境建議將 Origin 限制為指定域名
package middleware

import (
	"github.com/gin-gonic/gin"
)

// CORS 跨域中間件
// 開發環境允許所有來源（*），生產環境應限制
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*") // 允許來源，生產環境應改為具體域名
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Content-Type")
		c.Header("Access-Control-Max-Age", "86400") // 預檢請求快取 24 小時

		// OPTIONS 預檢請求直接返回 204
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
