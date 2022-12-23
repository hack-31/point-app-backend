package domain

import (
	"context"
	"time"

	"github.com/hack-31/point-app-backend/entity"
	"github.com/hack-31/point-app-backend/repository"
)

// Userに対するインターフェース
type UserRepo interface {
	FindUserByEmail(ctx context.Context, db repository.Queryer, e *string) (entity.User, error)
	RegisterUser(ctx context.Context, db repository.Execer, u *entity.User) error
	UpdatePassword(ctx context.Context, db repository.Execer, email, pass *string) error
}

// トークンに対するインターフェース
type TokenGenerator interface {
	GenerateToken(ctx context.Context, u entity.User) ([]byte, error)
}

// キャッシュに対するインターフェース
type Cache interface {
	Save(ctx context.Context, key, value string, minute time.Duration) error
	Load(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
}
