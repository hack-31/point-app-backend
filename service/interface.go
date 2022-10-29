package service

import (
	"context"

	"github.com/hack-31/point-app-backend/entity"
	"github.com/hack-31/point-app-backend/repository"
)

type UserRegister interface {
	RegisterUser(ctx context.Context, db repository.Execer, u *entity.User) error
}
