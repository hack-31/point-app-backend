package handler

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/hack-31/point-app-backend/domain/model"
	"github.com/hack-31/point-app-backend/service"
)

type RegisterUserService interface {
	RegisterUser(ctx context.Context, temporaryUserId, confirmCode string) (*model.User, string, error)
}

type RegisterTemporaryUserService interface {
	RegisterTemporaryUser(ctx context.Context, firstName, firstNameKana, familyName, familyNameKana, email, password string) (string, error)
}

type DeleteUserService interface {
	DeleteUser(ctx *gin.Context, userID model.UserID) error
}

type SigninService interface {
	Signin(ctx context.Context, email, password string) (string, error)
}

type GetUsersService interface {
	GetUsers(ctx context.Context) (service.GetUsersResponse, error)
}

type GetAccountService interface {
	GetAccount(ctx *gin.Context) (service.GetAccountResponse, error)
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

type UpdatePasswordService interface {
	UpdatePassword(ctx *gin.Context, oldPassword, newPassword string) error
}

type UpdateAccountService interface {
	UpdateAccount(ctx *gin.Context, familyName, familyNameKana, firstName, firstNameKana string) error
}

type UpdateTemporaryEmailService interface {
	UpdateTemporaryEmail(ctx *gin.Context, email string) (string, error)
}

type UpdateEmailService interface {
	UpdateEmail(ctx *gin.Context, temporaryEmailID, confirmCode string) error
}

type GetNotificationService interface {
	GetNotification(ctx *gin.Context, notificationID model.NotificationID) (service.GetNotificationResponse, error)
}

type GetNotificationsService interface {
	GetNotifications(ctx *gin.Context, nextToken, size string) (service.GetNotificationsResponse, error)
}

type GetUncheckedNotificationCountService interface {
	GetUncheckedNotificationCount(ctx *gin.Context) (<-chan int, error)
}
