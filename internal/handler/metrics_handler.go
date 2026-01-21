// 監控指標 HTTP 處理器
//
// 本檔案提供系統監控指標查詢功能
// 包含：資料庫連線池狀態、Redis 連線狀態
// 可擴展為 Prometheus 格式
package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/Mag1cFall/magtrade/internal/cache"
	"github.com/Mag1cFall/magtrade/internal/database"
	"github.com/gin-gonic/gin"
)

// MetricsHandler 監控指標處理器
type MetricsHandler struct{}

func NewMetricsHandler() *MetricsHandler {
	return &MetricsHandler{}
}

// GetMetrics 取得系統監控指標
// GET /metrics
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

// HealthCheck 完整健康檢查
// GET /health
func (h *MetricsHandler) HealthCheck(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	status := "ok"
	httpStatus := http.StatusOK

	// 資料庫檢查
	dbStatus := "ok"
	db := database.Get()
	sqlDB, err := db.DB()
	if err != nil {
		dbStatus = "error: " + err.Error()
		status = "degraded"
	} else if err := sqlDB.PingContext(ctx); err != nil {
		dbStatus = "error: " + err.Error()
		status = "degraded"
	}

	// Redis 檢查
	redisStatus := "ok"
	rdb := cache.Get()
	if rdb == nil {
		redisStatus = "not initialized"
		status = "degraded"
	} else if _, err := rdb.Ping(ctx).Result(); err != nil {
		redisStatus = "error: " + err.Error()
		status = "degraded"
	}

	if status == "degraded" {
		httpStatus = http.StatusServiceUnavailable
	}

	c.JSON(httpStatus, gin.H{
		"status": status,
		"components": gin.H{
			"database": dbStatus,
			"redis":    redisStatus,
		},
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

