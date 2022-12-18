package auth

import (
	"context"
	"time"
)

type Store interface {
	Save(ctx context.Context, key, value string, minute time.Duration) error
	Load(ctx context.Context, key string) (string, error)
	Expire(ctx context.Context, key string, expiration time.Duration) error
}
