package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	App      AppConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

// AppConfig holds application-level config
type AppConfig struct {
	Env  string
	Port string
}

// DatabaseConfig holds database connection config
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

// JWTConfig holds JWT config
type JWTConfig struct {
	Secret           string
	ExpiresIn        time.Duration
	RefreshExpiresIn time.Duration
}

// Load reads configuration from environment (and optionally .env file)
func Load() (*Config, error) {
	// Load .env file if it exists; ignore error if not present
	_ = godotenv.Load()

	expiresInHours, _ := strconv.Atoi(getEnv("JWT_EXPIRES_IN", "24"))
	refreshExpiresInHours, _ := strconv.Atoi(getEnv("JWT_REFRESH_EXPIRES_IN", "168"))

	cfg := &Config{
		App: AppConfig{
			Env:  getEnv("APP_ENV", "development"),
			Port: getEnv("APP_PORT", "8080"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			Name:     getEnv("DB_NAME", "todoapp"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:           getEnv("JWT_SECRET", "super-secret-key"),
			ExpiresIn:        time.Duration(expiresInHours) * time.Hour,
			RefreshExpiresIn: time.Duration(refreshExpiresInHours) * time.Hour,
		},
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
