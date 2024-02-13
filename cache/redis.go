package cache

import (
	"context"
	"fmt"
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

func (r Redis) SetUserInComputerGame(ctx context.Context, userID string, gameID string) error {
	_, err := r.rdb.Set(ctx, fmt.Sprintf("ingame:%s", userID), gameID, 0).Result()
	return err
}

func (r Redis) GetUserInComputerGame(ctx context.Context, userID string) (string, error) {
	return r.rdb.Get(ctx, fmt.Sprintf("ingame:%s", userID)).Result()
}

func (r Redis) DelUserInComputerGame(ctx context.Context, userID string) error {
	return r.rdb.Del(ctx, fmt.Sprintf("ingame:%s", userID)).Err()
}

func (r Redis) SetComputerGamestate(ctx context.Context, gameID string) error {
	_, err := r.rdb.RPush(ctx, fmt.Sprintf("gamestate:%s", gameID)).Result()
	return err
}

func (r Redis) GetComputerGamestate(ctx context.Context, gameID string) ([]string, error) {
	return r.rdb.LRange(ctx, fmt.Sprintf("gamestate:%s", gameID), 0, -1).Result()
}

func (r Redis) PushComputerGamestateMove(ctx context.Context, gameID, move string) error {
	_, err := r.rdb.RPush(ctx, fmt.Sprintf("gamestate:%s", gameID), move).Result()
	return err
}

func (r Redis) DelGamestate(ctx context.Context, gameID string) error {
	return r.rdb.Del(ctx, fmt.Sprintf("gamestate:%s", gameID)).Err()
}
