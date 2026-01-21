// 系統配置管理
//
// 本檔案使用 Viper 載入 YAML 配置檔
// 根據 APP_ENV 環境變數自動切換 dev/prod 配置
// 支援環境變數覆蓋（如 ${DB_PASSWORD}）
package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config 系統配置根結構
type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	Database  DatabaseConfig  `mapstructure:"database"`
	Redis     RedisConfig     `mapstructure:"redis"`
	Kafka     KafkaConfig     `mapstructure:"kafka"`
	JWT       JWTConfig       `mapstructure:"jwt"`
	AI        AIConfig        `mapstructure:"ai"`
	Log       LogConfig       `mapstructure:"log"`
	RateLimit RateLimitConfig `mapstructure:"rate_limit"`
	Email     EmailConfig     `mapstructure:"email"`
}

// ServerConfig HTTP 伺服器配置
type ServerConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	GRPCPort int    `mapstructure:"grpc_port"`
	Mode     string `mapstructure:"mode"` // debug/release
}

// DatabaseConfig PostgreSQL 資料庫配置
type DatabaseConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	DBName          string        `mapstructure:"dbname"`
	SSLMode         string        `mapstructure:"sslmode"` // disable/require
	Timezone        string        `mapstructure:"timezone"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`    // 閒置連線數
	MaxOpenConns    int           `mapstructure:"max_open_conns"`    // 最大連線數
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"` // 連線最大存活時間
}

// DSN 生成 PostgreSQL 連線字串
func (d *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		d.Host, d.Port, d.User, d.Password, d.DBName, d.SSLMode, d.Timezone,
	)
}

// RedisConfig Redis 快取配置
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

// Addr 生成 Redis 連線地址
func (r *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

// KafkaConfig Kafka 訊息佇列配置
type KafkaConfig struct {
	Brokers       []string          `mapstructure:"brokers"`
	ConsumerGroup string            `mapstructure:"consumer_group"`
	Topics        KafkaTopicsConfig `mapstructure:"topics"`
}

// KafkaTopicsConfig Kafka Topic 名稱配置
type KafkaTopicsConfig struct {
	FlashSaleOrders   string `mapstructure:"flash_sale_orders"`
	OrderStatusChange string `mapstructure:"order_status_change"`
	AIAnalysisTasks   string `mapstructure:"ai_analysis_tasks"`
}

// JWTConfig JWT 認證配置
type JWTConfig struct {
	Secret             string        `mapstructure:"secret"`
	AccessTokenExpire  time.Duration `mapstructure:"access_token_expire"`  // Access Token 有效期
	RefreshTokenExpire time.Duration `mapstructure:"refresh_token_expire"` // Refresh Token 有效期
}

// AIConfig AI 服務配置（SiliconFlow/OpenAI 兼容）
type AIConfig struct {
	Provider    string  `mapstructure:"provider"` // siliconflow/openai
	BaseURL     string  `mapstructure:"base_url"`
	APIKey      string  `mapstructure:"api_key"`
	Model       string  `mapstructure:"model"` // 模型名稱
	MaxTokens   int     `mapstructure:"max_tokens"`
	Temperature float64 `mapstructure:"temperature"` // 0-1，越高越隨機
}

// LogConfig 日誌配置
type LogConfig struct {
	Level  string `mapstructure:"level"`  // debug/info/warn/error
	Format string `mapstructure:"format"` // json/console
	Output string `mapstructure:"output"` // stdout/stderr/file
}

// RateLimitConfig 限流配置（令牌桶算法）
type RateLimitConfig struct {
	RequestsPerSecond int `mapstructure:"requests_per_second"` // 每秒請求數上限
	Burst             int `mapstructure:"burst"`               // 突發請求容量
}

// EmailConfig 郵件服務配置
type EmailConfig struct {
	SMTPHost     string `mapstructure:"smtp_host"`
	SMTPPort     int    `mapstructure:"smtp_port"`
	SMTPUser     string `mapstructure:"smtp_user"`
	SMTPPassword string `mapstructure:"smtp_password"`
	FromAddress  string `mapstructure:"from_address"`
	FromName     string `mapstructure:"from_name"`
}

var cfg *Config // 全域配置實例

// Load 載入配置檔
// 優先讀取 APP_ENV 環境變數選擇配置檔（預設 dev）
// 配置檔路徑優先順序：configPath → ./config → ../config → .
func Load(configPath string) (*Config, error) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}

	v := viper.New()
	v.SetConfigName(fmt.Sprintf("config.%s", env)) // config.dev.yaml 或 config.prod.yaml
	v.SetConfigType("yaml")

	if configPath != "" {
		v.AddConfigPath(configPath)
	}
	v.AddConfigPath("./config")
	v.AddConfigPath("../config")
	v.AddConfigPath(".")

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // 環境變數用底線替代點
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	cfg = &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	expandEnvVars(cfg) // 展開 ${ENV_VAR} 格式的環境變數

	return cfg, nil
}

// Get 取得配置實例（必須先呼叫 Load）
func Get() *Config {
	if cfg == nil {
		panic("config not loaded, call Load() first")
	}
	return cfg
}

// expandEnvVars 展開配置中的環境變數引用
// 支援兩種方式：直接環境變數（DB_HOST）或 YAML 中的 ${DB_HOST}
func expandEnvVars(c *Config) {
	// 資料庫配置
	if v := os.Getenv("DB_HOST"); v != "" {
		c.Database.Host = v
	} else {
		c.Database.Host = expandEnv(c.Database.Host)
	}
	if v := os.Getenv("DB_USER"); v != "" {
		c.Database.User = v
	} else {
		c.Database.User = expandEnv(c.Database.User)
	}
	if v := os.Getenv("DB_PASSWORD"); v != "" {
		c.Database.Password = v
	} else {
		c.Database.Password = expandEnv(c.Database.Password)
	}
	if v := os.Getenv("DB_NAME"); v != "" {
		c.Database.DBName = v
	} else {
		c.Database.DBName = expandEnv(c.Database.DBName)
	}

	// Redis 配置
	if v := os.Getenv("REDIS_HOST"); v != "" {
		c.Redis.Host = v
	} else {
		c.Redis.Host = expandEnv(c.Redis.Host)
	}
	if v := os.Getenv("REDIS_PASSWORD"); v != "" {
		c.Redis.Password = v
	} else {
		c.Redis.Password = expandEnv(c.Redis.Password)
	}

	// JWT 配置
	if v := os.Getenv("JWT_SECRET"); v != "" {
		c.JWT.Secret = v
	} else {
		c.JWT.Secret = expandEnv(c.JWT.Secret)
	}

	// AI 配置
	if v := os.Getenv("AI_API_KEY"); v != "" {
		c.AI.APIKey = v
	} else {
		c.AI.APIKey = expandEnv(c.AI.APIKey)
	}

	// Kafka 配置
	if v := os.Getenv("KAFKA_BROKER_1"); v != "" {
		c.Kafka.Brokers = []string{v}
		if v2 := os.Getenv("KAFKA_BROKER_2"); v2 != "" {
			c.Kafka.Brokers = append(c.Kafka.Brokers, v2)
		}
	} else {
		for i, broker := range c.Kafka.Brokers {
			c.Kafka.Brokers[i] = expandEnv(broker)
		}
	}

	// Email 配置
	if v := os.Getenv("SMTP_HOST"); v != "" {
		c.Email.SMTPHost = v
	} else {
		c.Email.SMTPHost = expandEnv(c.Email.SMTPHost)
	}
	if v := os.Getenv("SMTP_USER"); v != "" {
		c.Email.SMTPUser = v
		c.Email.FromAddress = v // 預設發件地址與使用者相同
	} else {
		c.Email.SMTPUser = expandEnv(c.Email.SMTPUser)
		c.Email.FromAddress = expandEnv(c.Email.FromAddress)
	}
	if v := os.Getenv("SMTP_PASSWORD"); v != "" {
		c.Email.SMTPPassword = v
	} else {
		c.Email.SMTPPassword = expandEnv(c.Email.SMTPPassword)
	}
}

// expandEnv 解析 ${VAR_NAME} 格式的環境變數引用
func expandEnv(s string) string {
	if strings.HasPrefix(s, "${") && strings.HasSuffix(s, "}") {
		envVar := s[2 : len(s)-1] // 提取變數名
		if val := os.Getenv(envVar); val != "" {
			return val
		}
	}
	return s
}
