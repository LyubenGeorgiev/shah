package cache

import (
	"context"
	"time"

	"github.com/LyubenGeorgiev/shah/db"
	"github.com/redis/go-redis/v9"
)

var gameFields = []string{"fen", "whiteid", "blackid"}

type Redis struct {
	rdb *redis.Client
}

func NewRedisCache() *Redis {
	host := db.Getenv("REDIS_HOST", "localhost")
	port := db.Getenv("REDIS_PORT", "6379")

	return &Redis{
		rdb: redis.NewClient(&redis.Options{
			Addr: host + ":" + port,
		}),
	}
}

func (r Redis) HealthCheck(ctx context.Context) error {
	return r.rdb.Ping(ctx).Err()
}

func (r Redis) Close() error {
	return r.rdb.Close()
}

func (r Redis) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	_, err := r.rdb.Set(ctx, key, value, expiration).Result()

	return err
}

func (r Redis) Del(ctx context.Context, key string) error {
	return r.rdb.Del(ctx, key).Err()
}

func (r Redis) Exists(ctx context.Context, key string, value string) bool {
	val, err := r.rdb.Get(ctx, key).Result()
	return err == nil && val == value
}
