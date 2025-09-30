package database

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Driver          string
	Host            string
	Port            string
	User            string
	Password        string
	Database        string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime int
}

func NewConfig() *Config {
	return &Config{
		Driver:          getEnv("DB_DRIVER", "postgres"),
		Host:            getEnv("DB_HOST", "localhost"),
		Port:            getEnv("DB_PORT", "5432"),
		User:            getEnv("DB_USER", "jiil"),
		Password:        getEnv("DB_PASSWORD", "juice"),
		Database:        getEnv("DB_NAME", "ecommerce"),
		SSLMode:         getEnv("DB_SSLMODE", "disable"),
		MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
		MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 5),
		ConnMaxLifetime: getEnvAsInt("DB_CONN_MAX_LIFETIME", 300),
	}
}

func (c *Config) GetDSN() string {
	switch c.Driver {
	case "postgres":
		return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			c.Host, c.Port, c.User, c.Password, c.Database, c.SSLMode)
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
			c.User, c.Password, c.Host, c.Port, c.Database)
	case "sqlite3":
		return c.Database
	default:
		return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			c.Host, c.Port, c.User, c.Password, c.Database, c.SSLMode)
	}
}

func (c *Config) GetDriver() string {
	return c.Driver
}

func (c *Config) GetMaxOpenConns() int {
	return c.MaxOpenConns
}

func (c *Config) GetMaxIdleConns() int {
	return c.MaxIdleConns
}

func (c *Config) GetConnMaxLifetime() int {
	return c.ConnMaxLifetime
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
