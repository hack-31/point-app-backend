package domain

import (
	"context"
	"time"

	"github.com/hack-31/point-app-backend/domain/model"
	"github.com/hack-31/point-app-backend/repository"
)

// Userに対するインターフェース
type UserRepo interface {
	FindUserByEmail(ctx context.Context, db repository.Queryer, e *string) (model.User, error)
	RegisterUser(ctx context.Context, db repository.Execer, u *model.User) error
	UpdatePassword(ctx context.Context, db repository.Execer, email, pass *string) error
	FindUsers(ctx context.Context, db repository.Queryer) (model.Users, error)
}

// ポイントに対するリポジトリインターフェース
type PointRepo interface {
	RegisterPointTransaction(ctx context.Context, db repository.Execer, fromUserID, toUserId model.UserID, sendPoint int) error
	UpdateSendablePoint(ctx context.Context, db repository.Execer, fromUserID model.UserID, sendPoint int) error
}

// トークンに対するインターフェース
type TokenGenerator interface {
	GenerateToken(ctx context.Context, u model.User) ([]byte, error)
}

// キャッシュに対するインターフェース
type Cache interface {
	Save(ctx context.Context, key, value string, minute time.Duration) error
	Load(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
}
