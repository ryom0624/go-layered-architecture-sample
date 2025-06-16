package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"layered-architecture-template/pkg/constants"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Auth     AuthConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	Driver   string
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type AuthConfig struct {
	JWTSecret            string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		// .env file is optional
	}

	config := &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", constants.DefaultServerPort),
		},
		Database: DatabaseConfig{
			Driver:   getEnv("DB_DRIVER", "postgres"),
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", constants.DefaultPostgreSQLPort),
			User:     getEnv("DB_USER", "user"),
			Password: getEnv("DB_PASSWORD", "password"),
			Name:     getEnv("DB_NAME", "database"),
		},
		Auth: AuthConfig{
			JWTSecret:            getEnv("JWT_SECRET", "your-secret-key"),
			AccessTokenDuration:  time.Duration(getEnvInt("ACCESS_TOKEN_DURATION_MINUTES", constants.DefaultAccessTokenDurationMinutes)) * time.Minute,
			RefreshTokenDuration: time.Duration(getEnvInt("REFRESH_TOKEN_DURATION_DAYS", constants.DefaultRefreshTokenDurationDays)) * constants.HoursPerDay * time.Hour,
		},
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}