package cache

import (
	"context"
	"time"
)

type Cache interface {
	HealthCheck(ctx context.Context) error
	Close() error
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
	Exists(ctx context.Context, key string, value string) bool
}
