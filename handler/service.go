package handler

import (
	"context"

	"github.com/hack-31/point-app-backend/entity"
)

type RegisterUserService interface {
	RegisterUser(ctx context.Context, temporaryUserId, confirmCode string) (*entity.User, string, error)
}

type RegisterTemporaryUserService interface {
	RegisterTemporaryUser(ctx context.Context, firstName, firstNameKana, familyName, familyNameKana, email, password string) (string, error)
}

type SigninService interface {
	Signin(ctx context.Context, email, password string) (string, error)
}

type GetUsersService interface {
	GetUsers(ctx context.Context) (entity.Users, error)
}

type SignoutService interface {
	Signout(ctx context.Context, userId string) error
}
