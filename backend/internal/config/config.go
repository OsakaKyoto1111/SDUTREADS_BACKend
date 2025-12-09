package config

import (
	"os"
)

type Config struct {
	AppPort     string
	DatabaseDSN string
	JWTSecret   string
	RedisAddr   string
}

func Load() (*Config, error) {
	return &Config{
		AppPort:     getEnv("APP_PORT", "8080"),
		DatabaseDSN: getEnv("DATABASE_DSN", "postgres://postgres:postgres@db:5432/backend?sslmode=disable"),
		JWTSecret:   getEnv("JWT_SECRET", "supersecret"),
		RedisAddr:   getEnv("REDIS_ADDR", "redis:6379"),
	}, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
