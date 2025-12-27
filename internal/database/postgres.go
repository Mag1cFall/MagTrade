// PostgreSQL 資料庫連線管理
//
// 本檔案負責資料庫初始化和連線池配置
// 包含：連線初始化、AutoMigrate、種子資料、連線池設定
// 使用 GORM 作為 ORM 框架
package database

import (
	"fmt"
	"os"
	"time"

	"github.com/Mag1cFall/magtrade/internal/config"
	"github.com/Mag1cFall/magtrade/internal/model"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

// Init 初始化資料庫連線
func Init(cfg *config.DatabaseConfig, log *zap.Logger) error {
	// 配置 GORM 日誌（慢查詢警告）
	gormLogger := logger.New(
		&zapWriter{log: log},
		logger.Config{
			SlowThreshold:             200 * time.Millisecond, // 慢查詢閾值
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	var err error
	db, err = gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
		Logger:                                   gormLogger,
		PrepareStmt:                              true, // 預編譯語句
		DisableForeignKeyConstraintWhenMigrating: false,
	})
	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}

	// 配置連線池
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// 測試連線
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Info("database connected",
		zap.String("host", cfg.Host),
		zap.Int("port", cfg.Port),
		zap.String("dbname", cfg.DBName),
	)

	return nil
}

// AutoMigrate 自動遷移資料表結構
func AutoMigrate() error {
	return db.AutoMigrate(
		&model.User{},
		&model.Product{},
		&model.FlashSale{},
		&model.Order{},
		&model.ChatHistory{},
		&model.AIRecommendation{},
	)
}

// SeedData 初始化種子資料
// 開發環境：自動建立 admin + 測試商品
// 生產環境：只建立 admin，密碼從環境變數讀取
func SeedData() error {
	var count int64
	db.Model(&model.User{}).Count(&count)
	if count > 0 { // 已有資料則跳過
		return nil
	}

	// 生產環境需設置 ADMIN_INIT_PASSWORD
	var password string
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "prod" || appEnv == "production" {
		password = os.Getenv("ADMIN_INIT_PASSWORD")
		if password == "" {
			return fmt.Errorf("ADMIN_INIT_PASSWORD is required in production")
		}
	} else {
		password = "mag1cfall1337"
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	admin := model.User{
		Username:      "admin",
		Email:         "admin@magtrade.com",
		PasswordHash:  string(hash),
		Role:          "admin",
		Status:        1,
		EmailVerified: true,
	}
	if err := db.Create(&admin).Error; err != nil {
		return err
	}

	// 開發環境建立測試商品
	if appEnv != "prod" && appEnv != "production" {
		products := []model.Product{
			{Name: "iPhone 15 Pro Max", Description: "Apple 旗舰智能手机", OriginalPrice: 9999.00, Status: model.ProductStatusOnShelf},
			{Name: "MacBook Pro 14", Description: "Apple M3 Pro 芯片笔记本电脑", OriginalPrice: 16999.00, Status: model.ProductStatusOnShelf},
			{Name: "Sony PS5", Description: "次世代游戏主机", OriginalPrice: 3999.00, Status: model.ProductStatusOnShelf},
		}
		for _, p := range products {
			db.Create(&p)
		}
	}

	return nil
}

// Get 取得資料庫連線（必須先呼叫 Init）
func Get() *gorm.DB {
	if db == nil {
		panic("database not initialized, call Init() first")
	}
	return db
}

// Close 關閉資料庫連線
func Close() error {
	if db == nil {
		return nil
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// zapWriter 將 GORM 日誌輸出到 Zap
type zapWriter struct {
	log *zap.Logger
}

func (w *zapWriter) Printf(format string, args ...interface{}) {
	w.log.Warn(fmt.Sprintf(format, args...))
}
