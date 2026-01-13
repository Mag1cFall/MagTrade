// MagTrade 高併發分散式秒殺系統 - 應用程式入口
//
// 本檔案負責初始化所有核心元件並啟動 HTTP 伺服器
// 啟動順序：配置 → 日誌 → 資料庫 → Redis → Kafka → WebSocket → HTTP
// 關閉時會優雅等待所有請求處理完畢（Graceful Shutdown）
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Mag1cFall/magtrade/internal/cache"
	"github.com/Mag1cFall/magtrade/internal/config"
	"github.com/Mag1cFall/magtrade/internal/database"
	magrpc "github.com/Mag1cFall/magtrade/internal/grpc"
	"github.com/Mag1cFall/magtrade/internal/handler"
	"github.com/Mag1cFall/magtrade/internal/logger"
	"github.com/Mag1cFall/magtrade/internal/mq"
	"github.com/Mag1cFall/magtrade/internal/pkg/utils"
	"github.com/Mag1cFall/magtrade/internal/router"
	"github.com/Mag1cFall/magtrade/internal/worker"
	"go.uber.org/zap"
)

func main() {
	// 載入配置檔，根據 APP_ENV 環境變數自動選擇 dev/prod
	cfg, err := config.Load("")
	if err != nil {
		fmt.Printf("failed to load config: %v\n", err)
		os.Exit(1)
	}

	// 初始化結構化日誌（Zap），支援 JSON 格式輸出
	if err := logger.Init(&cfg.Log); err != nil {
		fmt.Printf("failed to init logger: %v\n", err)
		os.Exit(1)
	}
	log := logger.Get()
	defer func() { _ = logger.Sync() }() // 確保程式結束前刷新日誌緩衝區

	log.Info("starting MagTrade server",
		zap.String("mode", cfg.Server.Mode),
		zap.Int("port", cfg.Server.Port),
	)

	// 初始化 PostgreSQL 連線池
	if err := database.Init(&cfg.Database, log); err != nil {
		log.Fatal("failed to init database", zap.Error(err))
	}
	defer database.Close()

	// GORM 自動遷移：根據 model 結構自動建立/更新表結構
	if err := database.AutoMigrate(); err != nil {
		log.Fatal("failed to auto migrate", zap.Error(err))
	}
	log.Info("database migration completed")

	// 插入種子資料（admin 帳號、測試商品），已存在則跳過
	if err := database.SeedData(); err != nil {
		log.Warn("failed to seed data", zap.Error(err))
	}

	// 初始化 Redis 連線，用於庫存快取和分散式鎖
	if err := cache.Init(&cfg.Redis, log); err != nil {
		log.Fatal("failed to init redis", zap.Error(err))
	}
	defer cache.Close()

	// 啟動 gRPC 庫存服務
	if cfg.Server.GRPCPort > 0 {
		grpcAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.GRPCPort)
		grpcServer, err := magrpc.StartGRPCServer(grpcAddr, log)
		if err != nil {
			log.Fatal("failed to start gRPC server", zap.Error(err))
		}
		defer grpcServer.GracefulStop()
	}

	// 初始化雪花演算法 ID 產生器，參數 1 是機器節點編號
	if err := utils.InitSnowflake(1); err != nil {
		log.Fatal("failed to init snowflake", zap.Error(err))
	}

	// Kafka 訊息佇列初始化
	producer := mq.NewProducer(&cfg.Kafka, log) // 生產者：發送秒殺訂單消息
	defer producer.Close()

	consumer := mq.NewConsumer(&cfg.Kafka, log) // 消費者：處理非同步訂單
	defer consumer.Close()

	// WebSocket Hub：管理所有 WS 連線，用於推送秒殺結果通知
	wsHub := handler.NewWSHub(log)
	go wsHub.Run() // 獨立 goroutine 執行訊息分發

	// 訂單 Worker：消費 Kafka 訊息，建立實際訂單並推送 WS 通知
	orderWorker := worker.NewOrderWorker(producer, wsHub, log)
	consumer.RegisterHandler(cfg.Kafka.Topics.FlashSaleOrders, orderWorker.HandleFlashSaleOrder)
	consumer.RegisterHandler(cfg.Kafka.Topics.OrderStatusChange, orderWorker.HandleOrderStatusChange)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	consumer.Start(ctx) // 啟動 Kafka 消費者

	// 定時任務 Worker：自動開啟/結束秒殺活動
	schedulerWorker := worker.NewSchedulerWorker(producer, log)
	schedulerWorker.Start(ctx)
	defer schedulerWorker.Stop()

	// 設定 Gin 路由
	r := router.Setup(cfg, producer, wsHub, log)

	// 配置 HTTP 伺服器
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  10 * time.Second, // 讀取請求超時
		WriteTimeout: 30 * time.Second, // 寫入回應超時（含秒殺處理時間）
		IdleTimeout:  60 * time.Second, // Keep-Alive 連線閒置超時
	}

	// 在獨立 goroutine 啟動伺服器，避免阻塞主流程
	go func() {
		log.Info("server starting", zap.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("failed to start server", zap.Error(err))
		}
	}()

	// 優雅關閉（Graceful Shutdown）
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 監聽 Ctrl+C 和 kill 信號
	<-quit                                               // 阻塞等待信號

	log.Info("shutting down server...")

	// 給 30 秒處理剩餘請求
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("server forced to shutdown", zap.Error(err))
	}

	cancel() // 取消 context，通知所有 goroutine 停止

	log.Info("server exited")
}
