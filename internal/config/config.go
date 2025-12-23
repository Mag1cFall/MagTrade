package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

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

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	DBName          string        `mapstructure:"dbname"`
	SSLMode         string        `mapstructure:"sslmode"`
	Timezone        string        `mapstructure:"timezone"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

func (d *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		d.Host, d.Port, d.User, d.Password, d.DBName, d.SSLMode, d.Timezone,
	)
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

func (r *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

type KafkaConfig struct {
	Brokers       []string          `mapstructure:"brokers"`
	ConsumerGroup string            `mapstructure:"consumer_group"`
	Topics        KafkaTopicsConfig `mapstructure:"topics"`
}

type KafkaTopicsConfig struct {
	FlashSaleOrders   string `mapstructure:"flash_sale_orders"`
	OrderStatusChange string `mapstructure:"order_status_change"`
	AIAnalysisTasks   string `mapstructure:"ai_analysis_tasks"`
}

type JWTConfig struct {
	Secret             string        `mapstructure:"secret"`
	AccessTokenExpire  time.Duration `mapstructure:"access_token_expire"`
	RefreshTokenExpire time.Duration `mapstructure:"refresh_token_expire"`
}

type AIConfig struct {
	Provider    string  `mapstructure:"provider"`
	BaseURL     string  `mapstructure:"base_url"`
	APIKey      string  `mapstructure:"api_key"`
	Model       string  `mapstructure:"model"`
	MaxTokens   int     `mapstructure:"max_tokens"`
	Temperature float64 `mapstructure:"temperature"`
}

type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
}

type RateLimitConfig struct {
	RequestsPerSecond int `mapstructure:"requests_per_second"`
	Burst             int `mapstructure:"burst"`
}

type EmailConfig struct {
	SMTPHost     string `mapstructure:"smtp_host"`
	SMTPPort     int    `mapstructure:"smtp_port"`
	SMTPUser     string `mapstructure:"smtp_user"`
	SMTPPassword string `mapstructure:"smtp_password"`
	FromAddress  string `mapstructure:"from_address"`
	FromName     string `mapstructure:"from_name"`
}

var cfg *Config

func Load(configPath string) (*Config, error) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}

	v := viper.New()
	v.SetConfigName(fmt.Sprintf("config.%s", env))
	v.SetConfigType("yaml")

	if configPath != "" {
		v.AddConfigPath(configPath)
	}
	v.AddConfigPath("./config")
	v.AddConfigPath("../config")
	v.AddConfigPath(".")

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	cfg = &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	expandEnvVars(cfg)

	return cfg, nil
}

func Get() *Config {
	if cfg == nil {
		panic("config not loaded, call Load() first")
	}
	return cfg
}

func expandEnvVars(c *Config) {
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
	if v := os.Getenv("JWT_SECRET"); v != "" {
		c.JWT.Secret = v
	} else {
		c.JWT.Secret = expandEnv(c.JWT.Secret)
	}
	if v := os.Getenv("AI_API_KEY"); v != "" {
		c.AI.APIKey = v
	} else {
		c.AI.APIKey = expandEnv(c.AI.APIKey)
	}
	if v := os.Getenv("KAFKA_BROKER_1"); v != "" {
		c.Kafka.Brokers = []string{v}
	} else {
		for i, broker := range c.Kafka.Brokers {
			c.Kafka.Brokers[i] = expandEnv(broker)
		}
	}
}

func expandEnv(s string) string {
	if strings.HasPrefix(s, "${") && strings.HasSuffix(s, "}") {
		envVar := s[2 : len(s)-1]
		if val := os.Getenv(envVar); val != "" {
			return val
		}
	}
	return s
}
