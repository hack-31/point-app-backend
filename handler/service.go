package handler

import (
	"context"

	"github.com/gin-gonic/gin"
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

type GetAccountService interface {
	GetAccount(ctx *gin.Context) (entity.User, error)
}

type SignoutService interface {
	Signout(ctx *gin.Context) error
}

type SendPointService interface {
	SendPoint(ctx *gin.Context, toUserId, sendPoint int) error
}

type ResetPasswordService interface {
	ResetPassword(ctx context.Context, email string) error
}
