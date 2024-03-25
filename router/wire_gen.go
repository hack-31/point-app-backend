// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package router

import (
	"github.com/hack-31/point-app-backend/domain"
	"github.com/hack-31/point-app-backend/handler"
	"github.com/hack-31/point-app-backend/repository"
	"github.com/hack-31/point-app-backend/service"
	"github.com/jmoiron/sqlx"
)

// Injectors from wire.go:

func InitHealthCheck() *handler.HealthCheck {
	healthCheck := handler.NewHealthCheckHandler()
	return healthCheck
}

func InitRegisterUser(db *sqlx.DB, rep domain.UserRepo, cache domain.Cache, jwter domain.TokenGenerator) *handler.RegisterUser {
	registerUser := service.NewRegisterUser(db, rep, cache, jwter)
	handlerRegisterUser := handler.NewRegisterUserHandler(registerUser)
	return handlerRegisterUser
}

func InitRegisterTemporaryUser(db *sqlx.DB, rep domain.UserRepo, cache domain.Cache) *handler.RegisterTemporaryUser {
	registerTemporaryUser := service.NewRegisterTemporaryUser(db, rep, cache)
	handlerRegisterTemporaryUser := handler.NewRegisterTemporaryUserHandler(registerTemporaryUser)
	return handlerRegisterTemporaryUser
}

func InitSignin(db *sqlx.DB, rep domain.UserRepo, cache domain.Cache, jwter domain.TokenGenerator) *handler.Signin {
	signin := service.NewSignin(db, rep, cache, jwter)
	handlerSignin := handler.NewSigninHandler(signin)
	return handlerSignin
}

func InitResetPassword(db *sqlx.DB, rep domain.UserRepo) *handler.ResetPassword {
	resetPassword := service.NewResetPassword(db, rep)
	handlerResetPassword := handler.NewResetPasswordHandler(resetPassword)
	return handlerResetPassword
}

func InitGetUsers(db *sqlx.DB, repo *repository.Repository, jwter domain.TokenGenerator) *handler.GetUsers {
	getUsers := service.NewGetUsers(db, repo, jwter)
	handlerGetUsers := handler.NewGetUsers(getUsers)
	return handlerGetUsers
}

func InitDeleteUser(transacter repository.Transacter, repo domain.UserRepo, cache domain.Cache) *handler.DeleteUser {
	deleteUser := service.NewDeleteUser(cache, repo, transacter)
	handlerDeleteUser := handler.NewDeleteUser(deleteUser)
	return handlerDeleteUser
}

func InitGetAccount(db *sqlx.DB, repo *repository.Repository) *handler.GetAccount {
	getAccount := service.NewGetAccount(db, repo)
	handlerGetAccount := handler.NewGetAccount(getAccount)
	return handlerGetAccount
}

func InitSignout(cache domain.Cache) *handler.Signout {
	signout := service.NewSignout(cache)
	handlerSignout := handler.NewSignoutHandler(signout)
	return handlerSignout
}

func InitSendPoint(repo *repository.Repository, connection *repository.AppConnection, cache domain.Cache) *handler.SendPoint {
	sendPoint := service.NewSendPoint(repo, connection, cache)
	handlerSendPoint := handler.NewSendPoint(sendPoint)
	return handlerSendPoint
}

func InitUpdatePassword(db *sqlx.DB, repo domain.UserRepo) *handler.UpdatePassword {
	updatePassword := service.NewUpdatePassword(db, repo)
	handlerUpdatePassword := handler.NewUpdatePasswordHandler(updatePassword)
	return handlerUpdatePassword
}

func InitUpdateAccount(db *sqlx.DB, repo domain.UserRepo) *handler.UpdateAccount {
	updateAccount := service.NewUpdateAccount(db, repo)
	handlerUpdateAccount := handler.NewUpdateAccountHandler(updateAccount)
	return handlerUpdateAccount
}

func InitUpdateTemporaryEmail(db *sqlx.DB, cache domain.Cache, repo domain.UserRepo) *handler.UpdateTemporaryEmail {
	updateTemporaryEmail := service.NewUpdateTemporaryEmail(db, cache, repo)
	handlerUpdateTemporaryEmail := handler.NewUpdateTemporaryEmailHandler(updateTemporaryEmail)
	return handlerUpdateTemporaryEmail
}

func InitUpdateEmail(db *sqlx.DB, cache domain.Cache, repo domain.UserRepo) *handler.UpdateEmail {
	updateEmail := service.NewUpdateEmail(db, cache, repo)
	handlerUpdateEmail := handler.NewUpdateEmailHandler(updateEmail)
	return handlerUpdateEmail
}

func InitGetNotification(cache domain.Cache, repo *repository.Repository, connection repository.Transacter) *handler.GetNotification {
	getNotification := service.NewGetNotification(cache, repo, connection)
	handlerGetNotification := handler.NewGetNotification(getNotification)
	return handlerGetNotification
}

func InitGetNotifications(db *sqlx.DB, repo *repository.Repository) *handler.GetNotifications {
	getNotifications := service.NewGetNotifications(db, repo)
	handlerGetNotifications := handler.NewGetNotifications(getNotifications)
	return handlerGetNotifications
}

func InitGetUncheckedNotificationCount(db *sqlx.DB, cache domain.Cache, repo domain.NotificationRepo) *handler.GetUncheckedNotificationCount {
	getUncheckedNotificationCount := service.NewGetUncheckedNotificationCount(db, cache, repo)
	handlerGetUncheckedNotificationCount := handler.NewGetUncheckedNotificationCount(getUncheckedNotificationCount)
	return handlerGetUncheckedNotificationCount
}
