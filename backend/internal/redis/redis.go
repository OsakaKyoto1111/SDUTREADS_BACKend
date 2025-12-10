package redis

import (
	"log"

	"backend/internal/config"
)

type Client struct{}

func InitRedis(cfg *config.Config) *Client {
	log.Printf("Redis stub initialized at %s", cfg.RedisAddr)
	return &Client{}
}
