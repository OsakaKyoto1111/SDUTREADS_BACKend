package redis

import (
    "log"

    "backend/internal/config"
)

// Client is a placeholder for a Redis client.
type Client struct{}

// InitRedis logs the configured address and returns a stub client.
func InitRedis(cfg *config.Config) *Client {
    log.Printf("Redis stub initialized at %s", cfg.RedisAddr)
    return &Client{}
}
