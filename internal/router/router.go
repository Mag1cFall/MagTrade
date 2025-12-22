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

func Setup(cfg *config.Config, producer *mq.Producer, wsHub *handler.WSHub, log *zap.Logger) *gin.Engine {
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(middleware.Recovery(log))
	r.Use(middleware.Logger(log))
	r.Use(middleware.CORS())
	r.Use(middleware.Security())
	r.Use(middleware.RequestSizeLimit(10 * 1024 * 1024))

	ipLimiter := middleware.NewIPRateLimiter(&cfg.RateLimit)
	r.Use(middleware.RateLimit(ipLimiter))

	authHandler := handler.NewAuthHandler(&cfg.JWT)
	productHandler := handler.NewProductHandler()
	anomalyDetector := ai.NewAnomalyDetector(log)
	flashSaleHandler := handler.NewFlashSaleHandler(producer, anomalyDetector, log)
	orderHandler := handler.NewOrderHandler(producer, log)
	aiHandler := handler.NewAIHandler(&cfg.AI, log)
	wsHandler := handler.NewWSHandler(wsHub, &cfg.JWT, log)

	captchaHandler := handler.NewCaptchaHandler()
	metricsHandler := handler.NewMetricsHandler()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"name":    "MagTrade API",
			"version": "1.0.0",
			"docs":    "/swagger",
		})
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.GET("/metrics", metricsHandler.GetMetrics)

	r.StaticFile("/swagger", "./docs/swagger.html")
	r.StaticFile("/swagger.yaml", "./docs/swagger.yaml")

	r.GET("/ws/notifications", wsHandler.HandleConnection)

	v1 := r.Group("/api/v1")
	{
		v1.GET("/captcha", captchaHandler.GetCaptcha)
		v1.GET("/captcha/check", captchaHandler.CheckNeedsCaptcha)

		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.Refresh)
			auth.GET("/me", middleware.Auth(&cfg.JWT), authHandler.Me)
		}

		products := v1.Group("/products")
		{
			products.GET("", productHandler.List)
			products.GET("/:id", productHandler.GetByID)
		}

		flashSales := v1.Group("/flash-sales")
		{
			flashSales.GET("", flashSaleHandler.List)
			flashSales.GET("/:id", flashSaleHandler.GetByID)
			flashSales.GET("/:id/stock", flashSaleHandler.GetStock)

			flashSaleRateLimiter := middleware.NewFlashSaleRateLimiter()
			flashSales.POST("/:id/rush",
				middleware.Auth(&cfg.JWT),
				middleware.FlashSaleRateLimit(flashSaleRateLimiter),
				flashSaleHandler.Rush,
			)
		}

		orders := v1.Group("/orders")
		orders.Use(middleware.Auth(&cfg.JWT))
		{
			orders.GET("", orderHandler.List)
			orders.GET("/:order_no", orderHandler.GetByOrderNo)
			orders.POST("/:order_no/pay", orderHandler.Pay)
			orders.POST("/:order_no/cancel", orderHandler.Cancel)
		}

		aiGroup := v1.Group("/ai")
		aiGroup.Use(middleware.Auth(&cfg.JWT))
		{
			aiGroup.POST("/chat", aiHandler.Chat)
			aiGroup.GET("/chat/history", aiHandler.GetChatHistory)
			aiGroup.DELETE("/chat/history", aiHandler.ClearChatHistory)
			aiGroup.GET("/recommendations/:flash_sale_id", aiHandler.GetRecommendation)
		}

		admin := v1.Group("/admin")
		admin.Use(middleware.Auth(&cfg.JWT), middleware.AdminOnly())
		{
			admin.POST("/products", productHandler.Create)
			admin.PUT("/products/:id", productHandler.Update)
			admin.DELETE("/products/:id", productHandler.Delete)

			admin.POST("/flash-sales", flashSaleHandler.Create)

			admin.POST("/ai/analyze/:flash_sale_id", aiHandler.TriggerAnalysis)
		}
	}

	return r
}
