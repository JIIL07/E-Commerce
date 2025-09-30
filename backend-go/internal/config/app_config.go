package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type AppConfig struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	Redis    RedisConfig    `json:"redis"`
	JWT      JWTConfig      `json:"jwt"`
	Stripe   StripeConfig   `json:"stripe"`
	Logging  LoggingConfig  `json:"logging"`
	Cache    CacheConfig    `json:"cache"`
	Metrics  MetricsConfig  `json:"metrics"`
}

type ServerConfig struct {
	Host         string        `json:"host"`
	Port         int           `json:"port"`
	ReadTimeout  time.Duration `json:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout"`
	IdleTimeout  time.Duration `json:"idle_timeout"`
	Environment  string        `json:"environment"`
}

type DatabaseConfig struct {
	Driver          string        `json:"driver"`
	Host            string        `json:"host"`
	Port            int           `json:"port"`
	User            string        `json:"user"`
	Password        string        `json:"password"`
	Database        string        `json:"database"`
	SSLMode         string        `json:"ssl_mode"`
	MaxOpenConns    int           `json:"max_open_conns"`
	MaxIdleConns    int           `json:"max_idle_conns"`
	ConnMaxLifetime time.Duration `json:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `json:"conn_max_idle_time"`
}

type RedisConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

type JWTConfig struct {
	Secret    string        `json:"secret"`
	ExpiresIn time.Duration `json:"expires_in"`
	RefreshIn time.Duration `json:"refresh_in"`
	Issuer    string        `json:"issuer"`
	Audience  string        `json:"audience"`
}

type StripeConfig struct {
	SecretKey      string `json:"secret_key"`
	WebhookSecret  string `json:"webhook_secret"`
	PublishableKey string `json:"publishable_key"`
}

type LoggingConfig struct {
	Level      string `json:"level"`
	Format     string `json:"format"`
	Output     string `json:"output"`
	Filename   string `json:"filename"`
	MaxSize    int    `json:"max_size"`
	MaxBackups int    `json:"max_backups"`
	MaxAge     int    `json:"max_age"`
	Compress   bool   `json:"compress"`
}

type CacheConfig struct {
	DefaultTTL      time.Duration `json:"default_ttl"`
	MaxSize         int           `json:"max_size"`
	CleanupInterval time.Duration `json:"cleanup_interval"`
}

type MetricsConfig struct {
	Enabled   bool   `json:"enabled"`
	Port      int    `json:"port"`
	Path      string `json:"path"`
	Namespace string `json:"namespace"`
}

var globalConfig *AppConfig

func LoadConfig(configPath string) (*AppConfig, error) {
	config := &AppConfig{}

	if configPath != "" {
		if err := loadFromFile(config, configPath); err != nil {
			return nil, fmt.Errorf("failed to load config from file: %w", err)
		}
	}

	loadFromEnv(config)
	setDefaults(config)

	globalConfig = config
	return config, nil
}

func GetConfig() *AppConfig {
	if globalConfig == nil {
		config, _ := LoadConfig("")
		return config
	}
	return globalConfig
}

func loadFromFile(config *AppConfig, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, config)
}

func loadFromEnv(config *AppConfig) {
	config.Server.Host = getEnv("SERVER_HOST", config.Server.Host)
	config.Server.Port = getEnvAsInt("SERVER_PORT", config.Server.Port)
	config.Server.Environment = getEnv("ENVIRONMENT", config.Server.Environment)

	config.Database.Driver = getEnv("DB_DRIVER", config.Database.Driver)
	config.Database.Host = getEnv("DB_HOST", config.Database.Host)
	config.Database.Port = getEnvAsInt("DB_PORT", config.Database.Port)
	config.Database.User = getEnv("DB_USER", config.Database.User)
	config.Database.Password = getEnv("DB_PASSWORD", config.Database.Password)
	config.Database.Database = getEnv("DB_NAME", config.Database.Database)
	config.Database.SSLMode = getEnv("DB_SSLMODE", config.Database.SSLMode)
	config.Database.MaxOpenConns = getEnvAsInt("DB_MAX_OPEN_CONNS", config.Database.MaxOpenConns)
	config.Database.MaxIdleConns = getEnvAsInt("DB_MAX_IDLE_CONNS", config.Database.MaxIdleConns)

	config.Redis.Host = getEnv("REDIS_HOST", config.Redis.Host)
	config.Redis.Port = getEnvAsInt("REDIS_PORT", config.Redis.Port)
	config.Redis.Password = getEnv("REDIS_PASSWORD", config.Redis.Password)
	config.Redis.DB = getEnvAsInt("REDIS_DB", config.Redis.DB)

	config.JWT.Secret = getEnv("JWT_SECRET", config.JWT.Secret)
	config.JWT.ExpiresIn = getEnvAsDuration("JWT_EXPIRES_IN", config.JWT.ExpiresIn)
	config.JWT.RefreshIn = getEnvAsDuration("JWT_REFRESH_IN", config.JWT.RefreshIn)
	config.JWT.Issuer = getEnv("JWT_ISSUER", config.JWT.Issuer)
	config.JWT.Audience = getEnv("JWT_AUDIENCE", config.JWT.Audience)

	config.Stripe.SecretKey = getEnv("STRIPE_SECRET_KEY", config.Stripe.SecretKey)
	config.Stripe.WebhookSecret = getEnv("STRIPE_WEBHOOK_SECRET", config.Stripe.WebhookSecret)
	config.Stripe.PublishableKey = getEnv("STRIPE_PUBLISHABLE_KEY", config.Stripe.PublishableKey)

	config.Logging.Level = getEnv("LOG_LEVEL", config.Logging.Level)
	config.Logging.Format = getEnv("LOG_FORMAT", config.Logging.Format)
	config.Logging.Output = getEnv("LOG_OUTPUT", config.Logging.Output)
	config.Logging.Filename = getEnv("LOG_FILENAME", config.Logging.Filename)
	config.Logging.MaxSize = getEnvAsInt("LOG_MAX_SIZE", config.Logging.MaxSize)
	config.Logging.MaxBackups = getEnvAsInt("LOG_MAX_BACKUPS", config.Logging.MaxBackups)
	config.Logging.MaxAge = getEnvAsInt("LOG_MAX_AGE", config.Logging.MaxAge)
	config.Logging.Compress = getEnvAsBool("LOG_COMPRESS", config.Logging.Compress)

	config.Cache.DefaultTTL = getEnvAsDuration("CACHE_DEFAULT_TTL", config.Cache.DefaultTTL)
	config.Cache.MaxSize = getEnvAsInt("CACHE_MAX_SIZE", config.Cache.MaxSize)
	config.Cache.CleanupInterval = getEnvAsDuration("CACHE_CLEANUP_INTERVAL", config.Cache.CleanupInterval)

	config.Metrics.Enabled = getEnvAsBool("METRICS_ENABLED", config.Metrics.Enabled)
	config.Metrics.Port = getEnvAsInt("METRICS_PORT", config.Metrics.Port)
	config.Metrics.Path = getEnv("METRICS_PATH", config.Metrics.Path)
	config.Metrics.Namespace = getEnv("METRICS_NAMESPACE", config.Metrics.Namespace)
}

func setDefaults(config *AppConfig) {
	if config.Server.Host == "" {
		config.Server.Host = "0.0.0.0"
	}
	if config.Server.Port == 0 {
		config.Server.Port = 5000
	}
	if config.Server.ReadTimeout == 0 {
		config.Server.ReadTimeout = 30 * time.Second
	}
	if config.Server.WriteTimeout == 0 {
		config.Server.WriteTimeout = 30 * time.Second
	}
	if config.Server.IdleTimeout == 0 {
		config.Server.IdleTimeout = 120 * time.Second
	}
	if config.Server.Environment == "" {
		config.Server.Environment = "development"
	}

	if config.Database.Driver == "" {
		config.Database.Driver = "postgres"
	}
	if config.Database.Host == "" {
		config.Database.Host = "localhost"
	}
	if config.Database.Port == 0 {
		config.Database.Port = 5432
	}
	if config.Database.User == "" {
		config.Database.User = "jiil"
	}
	if config.Database.Password == "" {
		config.Database.Password = "juice"
	}
	if config.Database.Database == "" {
		config.Database.Database = "ecommerce"
	}
	if config.Database.SSLMode == "" {
		config.Database.SSLMode = "disable"
	}
	if config.Database.MaxOpenConns == 0 {
		config.Database.MaxOpenConns = 25
	}
	if config.Database.MaxIdleConns == 0 {
		config.Database.MaxIdleConns = 5
	}
	if config.Database.ConnMaxLifetime == 0 {
		config.Database.ConnMaxLifetime = 5 * time.Minute
	}
	if config.Database.ConnMaxIdleTime == 0 {
		config.Database.ConnMaxIdleTime = 1 * time.Minute
	}

	if config.Redis.Host == "" {
		config.Redis.Host = "localhost"
	}
	if config.Redis.Port == 0 {
		config.Redis.Port = 6379
	}

	if config.JWT.Secret == "" {
		config.JWT.Secret = "your-secret-key"
	}
	if config.JWT.ExpiresIn == 0 {
		config.JWT.ExpiresIn = 24 * time.Hour
	}
	if config.JWT.RefreshIn == 0 {
		config.JWT.RefreshIn = 7 * 24 * time.Hour
	}
	if config.JWT.Issuer == "" {
		config.JWT.Issuer = "ecommerce-api"
	}
	if config.JWT.Audience == "" {
		config.JWT.Audience = "ecommerce-client"
	}

	if config.Logging.Level == "" {
		config.Logging.Level = "info"
	}
	if config.Logging.Format == "" {
		config.Logging.Format = "json"
	}
	if config.Logging.Output == "" {
		config.Logging.Output = "stdout"
	}

	if config.Cache.DefaultTTL == 0 {
		config.Cache.DefaultTTL = 1 * time.Hour
	}
	if config.Cache.MaxSize == 0 {
		config.Cache.MaxSize = 1000
	}
	if config.Cache.CleanupInterval == 0 {
		config.Cache.CleanupInterval = 10 * time.Minute
	}

	if config.Metrics.Port == 0 {
		config.Metrics.Port = 9090
	}
	if config.Metrics.Path == "" {
		config.Metrics.Path = "/metrics"
	}
	if config.Metrics.Namespace == "" {
		config.Metrics.Namespace = "ecommerce"
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func (c *AppConfig) IsDevelopment() bool {
	return strings.ToLower(c.Server.Environment) == "development"
}

func (c *AppConfig) IsProduction() bool {
	return strings.ToLower(c.Server.Environment) == "production"
}

func (c *AppConfig) IsTesting() bool {
	return strings.ToLower(c.Server.Environment) == "testing"
}

func (c *AppConfig) GetServerAddress() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

func (c *AppConfig) GetDatabaseDSN() string {
	switch c.Database.Driver {
	case "postgres":
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			c.Database.Host, c.Database.Port, c.Database.User, c.Database.Password, c.Database.Database, c.Database.SSLMode)
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
			c.Database.User, c.Database.Password, c.Database.Host, c.Database.Port, c.Database.Database)
	case "sqlite3":
		return c.Database.Database
	default:
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			c.Database.Host, c.Database.Port, c.Database.User, c.Database.Password, c.Database.Database, c.Database.SSLMode)
	}
}

func (c *AppConfig) GetRedisAddress() string {
	return fmt.Sprintf("%s:%d", c.Redis.Host, c.Redis.Port)
}
