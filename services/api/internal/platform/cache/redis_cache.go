package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

// ErrCacheMiss is returned when a key is not found in the cache.
var ErrCacheMiss = errors.New("cache: miss")

// RedisCache implements Cache using Redis as the backing store.
type RedisCache struct {
	client *goredis.Client
}

// NewRedisCache creates a RedisCache backed by the given Redis client.
func NewRedisCache(client *goredis.Client) *RedisCache {
	return &RedisCache{client: client}
}

// Get retrieves a value from Redis and unmarshals it into dest.
// Returns ErrCacheMiss if the key does not exist.
func (c *RedisCache) Get(ctx context.Context, key string, dest any) error {
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return ErrCacheMiss
		}
		return fmt.Errorf("cache: get: %w", err)
	}
	if err := json.Unmarshal([]byte(val), dest); err != nil {
		return fmt.Errorf("cache: unmarshal: %w", err)
	}
	return nil
}

// Set marshals value to JSON and stores it in Redis with the given TTL.
func (c *RedisCache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("cache: marshal: %w", err)
	}
	if err := c.client.Set(ctx, key, data, ttl).Err(); err != nil {
		return fmt.Errorf("cache: set: %w", err)
	}
	return nil
}

// Delete removes a key from Redis.
func (c *RedisCache) Delete(ctx context.Context, key string) error {
	if err := c.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("cache: delete: %w", err)
	}
	return nil
}
