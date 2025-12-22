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
	"github.com/Mag1cFall/magtrade/internal/handler"
	"github.com/Mag1cFall/magtrade/internal/logger"
	"github.com/Mag1cFall/magtrade/internal/mq"
	"github.com/Mag1cFall/magtrade/internal/pkg/utils"
	"github.com/Mag1cFall/magtrade/internal/router"
	"github.com/Mag1cFall/magtrade/internal/worker"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.Load("")
	if err != nil {
		fmt.Printf("failed to load config: %v\n", err)
		os.Exit(1)
	}

	if err := logger.Init(&cfg.Log); err != nil {
		fmt.Printf("failed to init logger: %v\n", err)
		os.Exit(1)
	}
	log := logger.Get()
	defer func() { _ = logger.Sync() }()

	log.Info("starting MagTrade server",
		zap.String("mode", cfg.Server.Mode),
		zap.Int("port", cfg.Server.Port),
	)

	if err := database.Init(&cfg.Database, log); err != nil {
		log.Fatal("failed to init database", zap.Error(err))
	}
	defer database.Close()

	if err := database.AutoMigrate(); err != nil {
		log.Fatal("failed to auto migrate", zap.Error(err))
	}
	log.Info("database migration completed")

	if err := cache.Init(&cfg.Redis, log); err != nil {
		log.Fatal("failed to init redis", zap.Error(err))
	}
	defer cache.Close()

	if err := utils.InitSnowflake(1); err != nil {
		log.Fatal("failed to init snowflake", zap.Error(err))
	}

	producer := mq.NewProducer(&cfg.Kafka, log)
	defer producer.Close()

	consumer := mq.NewConsumer(&cfg.Kafka, log)
	defer consumer.Close()

	wsHub := handler.NewWSHub(log)
	go wsHub.Run()

	orderWorker := worker.NewOrderWorker(producer, wsHub, log)
	consumer.RegisterHandler(cfg.Kafka.Topics.FlashSaleOrders, orderWorker.HandleFlashSaleOrder)
	consumer.RegisterHandler(cfg.Kafka.Topics.OrderStatusChange, orderWorker.HandleOrderStatusChange)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	consumer.Start(ctx)

	schedulerWorker := worker.NewSchedulerWorker(producer, log)
	schedulerWorker.Start(ctx)
	defer schedulerWorker.Stop()

	r := router.Setup(cfg, producer, wsHub, log)

	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Info("server starting", zap.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("failed to start server", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("server forced to shutdown", zap.Error(err))
	}

	cancel()

	log.Info("server exited")
}
