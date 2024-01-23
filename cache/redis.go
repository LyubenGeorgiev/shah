package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/LyubenGeorgiev/shah/cache/models"
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

func (r Redis) Exists(ctx context.Context, key string, value string) bool {
	val, err := r.rdb.Get(ctx, key).Result()
	return err == nil && val == value
}

func (r Redis) UserIsIngame(ctx context.Context, userID string) (string, error) {
	return r.rdb.HGet(ctx, userID, "gameid").Result()
}

func (r Redis) GetGame(ctx context.Context, gameID string) (*models.Game, error) {
	gameData, err := r.rdb.HMGet(ctx, hashesKeyFromGameID(gameID), gameFields...).Result()
	if err != nil {
		return nil, err
	}

	return gameFromSlice(gameData)
}

func hashesKeyFromGameID(gameID string) string {
	return fmt.Sprintf("game:%s", gameID)
}

func gameFromSlice(gameData []interface{}) (*models.Game, error) {
	BoardFEN, ok := gameData[0].(string)
	if !ok {
		return nil, fmt.Errorf("Missing fen field in redis for given game")
	}
	whiteUserID, ok := gameData[1].(string)
	if !ok {
		return nil, fmt.Errorf("Missing whiteID field in redis for given game")
	}
	blackUserID, ok := gameData[2].(string)
	if !ok {
		return nil, fmt.Errorf("Missing blackID field in redis for given game")
	}

	return &models.Game{
		BoardFEN:    BoardFEN,
		WhiteUserID: whiteUserID,
		BlackUserID: blackUserID,
	}, nil
}
