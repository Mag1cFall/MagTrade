// HTTP 路由配置
//
// 本檔案配置 Gin 路由和中間件
// 包含：認證、商品、秒殺、訂單、AI、管理員等 API 分組
// 公開路由、認證路由、管理員路由三層權限
package router

import (
	"net/http"

	"github.com/Mag1cFall/magtrade/internal/config"
	"github.com/Mag1cFall/magtrade/internal/handler"
	"github.com/Mag1cFall/magtrade/internal/middleware"
	"github.com/Mag1cFall/magtrade/internal/mq"
	"github.com/Mag1cFall/magtrade/internal/service/ai"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Setup 配置並返回 Gin 引擎
func Setup(cfg *config.Config, producer *mq.Producer, wsHub *handler.WSHub, log *zap.Logger) *gin.Engine {
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// 全局中間件
	r.Use(middleware.Recovery(log))                      // Panic 恢復
	r.Use(middleware.Logger(log))                        // 請求日誌
	r.Use(middleware.CORS())                             // 跨域處理
	r.Use(middleware.Security())                         // 安全標頭
	r.Use(middleware.RequestSizeLimit(10 * 1024 * 1024)) // 請求體限制 10MB

	// IP 限流中間件
	ipLimiter := middleware.NewIPRateLimiter(&cfg.RateLimit)
	r.Use(middleware.RateLimit(ipLimiter))

	// 初始化 Handler
	authHandler := handler.NewAuthHandler(&cfg.JWT, &cfg.Email)
	productHandler := handler.NewProductHandler()
	anomalyDetector := ai.NewAnomalyDetector(log)
	flashSaleHandler := handler.NewFlashSaleHandler(producer, anomalyDetector, log)
	orderHandler := handler.NewOrderHandler(producer, log)
	aiHandler := handler.NewAIHandler(&cfg.AI, log)
	wsHandler := handler.NewWSHandler(wsHub, &cfg.JWT, log)

	captchaHandler := handler.NewCaptchaHandler()
	metricsHandler := handler.NewMetricsHandler()
	uploadHandler := handler.NewUploadHandler()

	// 靜態檔案
	r.Static("/uploads", "./uploads")

	// 根路徑
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"name":    "MagTrade API",
			"version": "1.0.0",
			"docs":    "/swagger",
		})
	})

	// 健康檢查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// 監控指標
	r.GET("/metrics", metricsHandler.GetMetrics)

	// Swagger 文件
	r.StaticFile("/swagger", "./docs/swagger.html")
	r.StaticFile("/swagger.yaml", "./docs/swagger.yaml")

	// WebSocket 通知
	r.GET("/ws/notifications", wsHandler.HandleConnection)

	// ===== API v1 路由分組 =====
	v1 := r.Group("/api/v1")
	{
		// 驗證碼（公開）
		v1.GET("/captcha", captchaHandler.GetCaptcha)
		v1.GET("/captcha/check", captchaHandler.CheckNeedsCaptcha)

		// 認證（公開）
		auth := v1.Group("/auth")
		{
			auth.POST("/send-code", authHandler.SendEmailCode)
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.Refresh)
			auth.GET("/me", middleware.Auth(&cfg.JWT), authHandler.Me)
		}

		// 商品（公開）
		products := v1.Group("/products")
		{
			products.GET("", productHandler.List)
			products.GET("/:id", productHandler.GetByID)
		}

		// 秒殺活動
		flashSales := v1.Group("/flash-sales")
		{
			flashSales.GET("", flashSaleHandler.List)
			flashSales.GET("/:id", flashSaleHandler.GetByID)
			flashSales.GET("/:id/stock", flashSaleHandler.GetStock)

			// 搶購需要認證和特殊限流
			flashSaleRateLimiter := middleware.NewFlashSaleRateLimiter()
			flashSales.POST("/:id/rush",
				middleware.Auth(&cfg.JWT),
				middleware.FlashSaleRateLimit(flashSaleRateLimiter),
				flashSaleHandler.Rush,
			)
		}

		// 訂單（需認證）
		orders := v1.Group("/orders")
		orders.Use(middleware.Auth(&cfg.JWT))
		{
			orders.GET("", orderHandler.List)
			orders.GET("/:order_no", orderHandler.GetByOrderNo)
			orders.POST("/:order_no/pay", orderHandler.Pay)
			orders.POST("/:order_no/cancel", orderHandler.Cancel)
		}

		// AI 功能（需認證）
		aiGroup := v1.Group("/ai")
		aiGroup.Use(middleware.Auth(&cfg.JWT))
		{
			aiGroup.POST("/chat", aiHandler.Chat)
			aiGroup.POST("/chat/stream", aiHandler.ChatStream)
			aiGroup.GET("/chat/history", aiHandler.GetChatHistory)
			aiGroup.DELETE("/chat/history", aiHandler.ClearChatHistory)
			aiGroup.GET("/recommendations/:flash_sale_id", aiHandler.GetRecommendation)
		}

		// 管理員專用（需認證 + admin 角色）
		admin := v1.Group("/admin")
		admin.Use(middleware.Auth(&cfg.JWT), middleware.AdminOnly())
		{
			admin.POST("/products", productHandler.Create)
			admin.PUT("/products/:id", productHandler.Update)
			admin.DELETE("/products/:id", productHandler.Delete)

			admin.POST("/upload", uploadHandler.Upload)

			admin.POST("/flash-sales", flashSaleHandler.Create)

			admin.POST("/ai/analyze/:flash_sale_id", aiHandler.TriggerAnalysis)
		}
	}

	return r
}
