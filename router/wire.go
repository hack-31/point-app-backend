//go:build wireinject
// +build wireinject

package router

import (
	"github.com/google/wire"
	"github.com/hack-31/point-app-backend/domain"
	"github.com/hack-31/point-app-backend/handler"
	"github.com/hack-31/point-app-backend/repository"
	"github.com/hack-31/point-app-backend/service"
	"github.com/jmoiron/sqlx"
)

func InitHealthCheck() *handler.HealthCheck {
	wire.Build(handler.NewHealthCheckHandler)
	return &handler.HealthCheck{}
}

func InitRegisterUser(db *sqlx.DB, rep domain.UserRepo, cache domain.Cache, jwter domain.TokenGenerator) *handler.RegisterUser {
	wire.Build(
		handler.NewRegisterUserHandler,
		service.NewRegisterUser,
		wire.Bind(new(handler.RegisterUserService), new(*service.RegisterUser)),
	)
	return &handler.RegisterUser{}
}

func InitRegisterTemporaryUser(db *sqlx.DB, rep domain.UserRepo, cache domain.Cache) *handler.RegisterTemporaryUser {
	wire.Build(
		handler.NewRegisterTemporaryUserHandler,
		service.NewRegisterTemporaryUser,
		wire.Bind(new(handler.RegisterTemporaryUserService), new(*service.RegisterTemporaryUser)),
	)
	return &handler.RegisterTemporaryUser{}
}

func InitSignin(db *sqlx.DB, rep domain.UserRepo, cache domain.Cache, jwter domain.TokenGenerator) *handler.Signin {
	wire.Build(
		handler.NewSigninHandler,
		service.NewSignin,
		wire.Bind(new(handler.SigninService), new(*service.Signin)),
	)
	return &handler.Signin{}
}

func InitResetPassword(db *sqlx.DB, rep domain.UserRepo) *handler.ResetPassword {
	wire.Build(
		handler.NewResetPasswordHandler,
		service.NewResetPassword,
		wire.Bind(new(handler.ResetPasswordService), new(*service.ResetPassword)),
	)
	return &handler.ResetPassword{}
}

func InitGetUsers(db *sqlx.DB, repo domain.UserRepo, jwter domain.TokenGenerator) *handler.GetUsers {
	wire.Build(
		handler.NewGetUsers,
		service.NewGetUsers,
		wire.Bind(new(handler.GetUsersService), new(*service.GetUsers)),
	)
	return &handler.GetUsers{}
}

func InitDeleteUser(transacter repository.Transacter, repo domain.UserRepo, cache domain.Cache) *handler.DeleteUser {
	wire.Build(
		handler.NewDeleteUser,
		service.NewDeleteUser,
		wire.Bind(new(handler.DeleteUserService), new(*service.DeleteUser)),
	)
	return &handler.DeleteUser{}
}

func InitGetAccount(db *sqlx.DB, repo domain.UserRepo) *handler.GetAccount {
	wire.Build(
		handler.NewGetAccount,
		service.NewGetAccount,
		wire.Bind(new(handler.GetAccountService), new(*service.GetAccount)),
	)
	return &handler.GetAccount{}
}

func InitSignout(cache domain.Cache) *handler.Signout {
	wire.Build(
		handler.NewSignoutHandler,
		service.NewSignout,
		wire.Bind(new(handler.SignoutService), new(*service.Signout)),
	)
	return &handler.Signout{}
}

func InitSendPoint(repo *repository.Repository, connection *repository.AppConnection, cache domain.Cache) *handler.SendPoint {
	wire.Build(
		handler.NewSendPoint,
		service.NewSendPoint,
		wire.Bind(new(handler.SendPointService), new(*service.SendPoint)),
	)
	return &handler.SendPoint{}
}

func InitUpdatePassword(db *sqlx.DB, repo domain.UserRepo) *handler.UpdatePassword {
	wire.Build(
		handler.NewUpdatePasswordHandler,
		service.NewUpdatePassword,
		wire.Bind(new(handler.UpdatePasswordService), new(*service.UpdatePassword)),
	)
	return &handler.UpdatePassword{}
}

func InitUpdateAccount(db *sqlx.DB, repo domain.UserRepo) *handler.UpdateAccount {
	wire.Build(
		handler.NewUpdateAccountHandler,
		service.NewUpdateAccount,
		wire.Bind(new(handler.UpdateAccountService), new(*service.UpdateAccount)),
	)
	return &handler.UpdateAccount{}
}

func InitUpdateTemporaryEmail(db *sqlx.DB, cache domain.Cache, repo domain.UserRepo) *handler.UpdateTemporaryEmail {
	wire.Build(
		handler.NewUpdateTemporaryEmailHandler,
		service.NewUpdateTemporaryEmail,
		wire.Bind(new(handler.UpdateTemporaryEmailService), new(*service.UpdateTemporaryEmail)),
	)
	return &handler.UpdateTemporaryEmail{}
}

func InitUpdateEmail(db *sqlx.DB, cache domain.Cache, repo domain.UserRepo) *handler.UpdateEmail {
	wire.Build(
		handler.NewUpdateEmailHandler,
		service.NewUpdateEmail,
		wire.Bind(new(handler.UpdateEmailService), new(*service.UpdateEmail)),
	)
	return &handler.UpdateEmail{}
}

func InitGetNotification(cache domain.Cache, repo *repository.Repository, connection repository.Transacter) *handler.GetNotification {
	wire.Build(
		handler.NewGetNotification,
		service.NewGetNotification,
		wire.Bind(new(handler.GetNotificationService), new(*service.GetNotification)),
	)
	return &handler.GetNotification{}
}

func InitGetNotifications(db *sqlx.DB, repo *repository.Repository) *handler.GetNotifications {
	wire.Build(
		handler.NewGetNotifications,
		service.NewGetNotifications,
		wire.Bind(new(handler.GetNotificationsService), new(*service.GetNotifications)),
	)
	return &handler.GetNotifications{}
}

func InitGetUncheckedNotificationCount(db *sqlx.DB, cache domain.Cache, repo domain.NotificationRepo) *handler.GetUncheckedNotificationCount {
	wire.Build(
		handler.NewGetUncheckedNotificationCount,
		service.NewGetUncheckedNotificationCount,
		wire.Bind(new(handler.GetUncheckedNotificationCountService), new(*service.GetUncheckedNotificationCount)),
	)
	return &handler.GetUncheckedNotificationCount{}
}
