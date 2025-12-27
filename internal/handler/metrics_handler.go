// 監控指標 HTTP 處理器
//
// 本檔案提供系統監控指標查詢功能
// 目前僅包含資料庫連線池狀態
// 可擴展為 Prometheus 格式
package handler

import (
	"net/http"

	"github.com/Mag1cFall/magtrade/internal/database"
	"github.com/gin-gonic/gin"
)

// MetricsHandler 監控指標處理器
type MetricsHandler struct{}

func NewMetricsHandler() *MetricsHandler {
	return &MetricsHandler{}
}

// GetMetrics 取得系統監控指標
// GET /api/v1/metrics
func (h *MetricsHandler) GetMetrics(c *gin.Context) {
	db := database.Get()
	sqlDB, err := db.DB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get database connection"})
		return
	}

	stats := sqlDB.Stats() // 資料庫連線池統計

	c.JSON(http.StatusOK, gin.H{
		"database": gin.H{
			"max_open_connections": stats.MaxOpenConnections, // 最大連線數
			"open_connections":     stats.OpenConnections,    // 當前開啟
			"in_use":               stats.InUse,              // 使用中
			"idle":                 stats.Idle,               // 閒置
			"wait_count":           stats.WaitCount,          // 等待次數
			"wait_duration_ms":     stats.WaitDuration.Milliseconds(),
			"max_idle_closed":      stats.MaxIdleClosed,     // 因超過 MaxIdleConns 而關閉
			"max_idle_time_closed": stats.MaxIdleTimeClosed, // 因閒置超時關閉
			"max_lifetime_closed":  stats.MaxLifetimeClosed, // 因存活超時關閉
		},
	})
}
