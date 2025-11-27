package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config holds settings loaded from environment variables.
type Config struct {
	AppPort     string
	DatabaseDSN string
	JWTSecret   string
	RedisAddr   string
}

// Load reads configuration from environment variables and applies defaults.
func Load() (*Config, error) {
	_ = godotenv.Load()

	return &Config{
		AppPort:     getEnv("APP_PORT", "8080"),
		DatabaseDSN: getEnv("DATABASE_DSN", "host=localhost user=postgres password=postgres dbname=backend sslmode=disable"),
		JWTSecret:   getEnv("JWT_SECRET", "supersecret"),
		RedisAddr:   getEnv("REDIS_ADDR", "localhost:6379"),
	}, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
