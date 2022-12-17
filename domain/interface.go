package domain

import (
	"context"

	"github.com/hack-31/point-app-backend/entity"
	"github.com/hack-31/point-app-backend/repository"
)

// エンドポイントごとにインターフェースを作る

type RegisterTemporaryUserRep interface {
	FindUserByEmail(ctx context.Context, db repository.Queryer, e *string) (entity.User, error)
}

type UserRegister interface {
	RegisterUser(ctx context.Context, db repository.Execer, u *entity.User) error
}
