package database

import (
	"fmt"
	"time"

	"github.com/Mag1cFall/magtrade/internal/config"
	"github.com/Mag1cFall/magtrade/internal/model"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func Init(cfg *config.DatabaseConfig, log *zap.Logger) error {
	gormLogger := logger.New(
		&zapWriter{log: log},
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	var err error
	db, err = gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
		Logger:                                   gormLogger,
		PrepareStmt:                              true,
		DisableForeignKeyConstraintWhenMigrating: false,
	})
	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

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

func SeedData() error {
	var count int64
	db.Model(&model.User{}).Count(&count)
	if count > 0 {
		return nil
	}

	admin := model.User{
		Username:      "admin",
		Email:         "admin@magtrade.com",
		PasswordHash:  "$2a$10$bKmO1qYASj5loeB5oT4nBO2NXISkf2E7vaOPQQ5/pBXy0FSQpQO0m",
		Role:          "admin",
		Status:        1,
		EmailVerified: true,
	}
	if err := db.Create(&admin).Error; err != nil {
		return err
	}

	products := []model.Product{
		{Name: "iPhone 15 Pro Max", Description: "Apple 旗舰智能手机", OriginalPrice: 9999.00, Status: model.ProductStatusOnShelf},
		{Name: "MacBook Pro 14", Description: "Apple M3 Pro 芯片笔记本电脑", OriginalPrice: 16999.00, Status: model.ProductStatusOnShelf},
		{Name: "Sony PS5", Description: "次世代游戏主机", OriginalPrice: 3999.00, Status: model.ProductStatusOnShelf},
	}
	for _, p := range products {
		db.Create(&p)
	}

	return nil
}

func Get() *gorm.DB {
	if db == nil {
		panic("database not initialized, call Init() first")
	}
	return db
}

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

type zapWriter struct {
	log *zap.Logger
}

func (w *zapWriter) Printf(format string, args ...interface{}) {
	w.log.Warn(fmt.Sprintf(format, args...))
}
