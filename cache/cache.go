package cache

import (
	"context"
	"time"
)

type Cache interface {
	HealthCheck(ctx context.Context) error
	Close() error
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
	Del(ctx context.Context, key string) error
	Exists(ctx context.Context, key string, value string) bool

	SetUserInComputerGame(ctx context.Context, userID string, gameID string) error
	GetUserInComputerGame(ctx context.Context, userID string) (string, error)
	DelUserInComputerGame(ctx context.Context, userID string) error

	SetComputerGamestate(ctx context.Context, gameID string) error
	GetComputerGamestate(ctx context.Context, gameID string) ([]string, error)
	PushComputerGamestateMove(ctx context.Context, gameID, move string) error
	DelGamestate(ctx context.Context, gameID string) error
}
