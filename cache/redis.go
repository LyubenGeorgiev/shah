package cache

import (
	"context"

	"github.com/LyubenGeorgiev/shah/db"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	rdb *redis.Client
}

func NewRedisCache() *Redis {
	host := db.Getenv("REDIS_HOST", "localhost")
	port := db.Getenv("REDIS_PORT", "6379")
	pass := db.Getenv("REDIS_PASSWORD", "eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81")

	return &Redis{
		rdb: redis.NewClient(&redis.Options{
			Addr:     host + ":" + port,
			Password: pass,
		}),
	}
}

func (r Redis) HealthCheck(ctx context.Context) error {
	return r.rdb.Ping(ctx).Err()
}

func (r Redis) Close() error {
	return r.rdb.Close()
}
