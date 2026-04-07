package cache

import (
	"context"
	"time"
)

// Cache defines a generic key-value cache interface.
type Cache interface {
	Get(ctx context.Context, key string, dest any) error
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}
