package cache

import "context"

type Cache interface {
	HealthCheck(ctx context.Context) error
	Close() error
}
