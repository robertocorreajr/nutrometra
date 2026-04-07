package redis

import (
	"context"
	"fmt"

	goredis "github.com/redis/go-redis/v9"
	"nutrometra/api/internal/platform/config"
)

// Client is an exported alias for goredis.Client.
type Client = goredis.Client

// New creates and validates a Redis client.
func New(ctx context.Context, cfg config.RedisConfig) (*Client, error) {
	client := goredis.NewClient(&goredis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
	})
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis: ping failed: %w", err)
	}
	return client, nil
}
