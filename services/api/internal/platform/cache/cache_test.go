package cache_test

import (
	"context"
	"os"
	"testing"
	"time"

	"nutrometra/api/internal/platform/cache"

	goredis "github.com/redis/go-redis/v9"
)

func redisClient(t *testing.T) *goredis.Client {
	t.Helper()
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}
	client := goredis.NewClient(&goredis.Options{Addr: addr})
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		t.Skipf("skipping: Redis not available at %s: %v", addr, err)
	}
	t.Cleanup(func() { client.Close() })
	return client
}

func TestRedisCache_GetMiss(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	client := redisClient(t)
	c := cache.NewRedisCache(client)

	var dest string
	err := c.Get(context.Background(), "nonexistent-key-"+t.Name(), &dest)
	if err == nil {
		t.Fatal("expected ErrCacheMiss, got nil")
	}
	if err != cache.ErrCacheMiss {
		t.Fatalf("expected ErrCacheMiss, got %v", err)
	}
}

func TestRedisCache_SetGetRoundtrip(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	client := redisClient(t)
	c := cache.NewRedisCache(client)
	ctx := context.Background()
	key := "test-roundtrip-" + t.Name()

	type sample struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	input := sample{Name: "alpha", Value: 42}
	if err := c.Set(ctx, key, input, 10*time.Second); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	t.Cleanup(func() { client.Del(ctx, key) })

	var output sample
	if err := c.Get(ctx, key, &output); err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if output.Name != input.Name || output.Value != input.Value {
		t.Fatalf("roundtrip mismatch: got %+v, want %+v", output, input)
	}
}

func TestRedisCache_Delete(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	client := redisClient(t)
	c := cache.NewRedisCache(client)
	ctx := context.Background()
	key := "test-delete-" + t.Name()

	if err := c.Set(ctx, key, "hello", 10*time.Second); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	if err := c.Delete(ctx, key); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	var dest string
	err := c.Get(ctx, key, &dest)
	if err != cache.ErrCacheMiss {
		t.Fatalf("expected ErrCacheMiss after Delete, got %v", err)
	}
}
