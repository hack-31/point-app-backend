package handler

import (
	"context"

	"github.com/hack-31/point-app-backend/entity"
)

type RegisterUserService interface {
	RegisterUser(ctx context.Context, temporaryUserId, confirmCode string) (*entity.User, error)
}

type RegisterTemporaryUserService interface {
	RegisterTemporaryUser(ctx context.Context, firstName, firstNameKana, familyName, familyNameKana, email, password string) (string, error)
}
