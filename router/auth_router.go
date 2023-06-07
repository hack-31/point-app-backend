package router

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/hack-31/point-app-backend/auth"
	"github.com/hack-31/point-app-backend/config"
	"github.com/hack-31/point-app-backend/handler"
	"github.com/hack-31/point-app-backend/repository"
	"github.com/hack-31/point-app-backend/service"
	"github.com/hack-31/point-app-backend/utils/clock"
	"github.com/jmoiron/sqlx"
)

// アクセストークンを必要とする
// 認証が必要なルーティングの設定を行う
//
// @param
// ctx コンテキスト
// router ルーター
func SetAuthRouting(ctx context.Context, db *sqlx.DB, router *gin.Engine, cfg *config.Config) error {
	// レポジトリ
	clocker := clock.RealClocker{}
	rep := repository.NewRepository(clocker)

	// 一時保存をするキャッシュ
	cache, err := repository.NewKVS(ctx, cfg, repository.TemporaryRegister)
	if err != nil {
		return err
	}

	// トークン保存をするキャッシュ
	tokenCache, err := repository.NewKVS(ctx, cfg, repository.JWT)
	if err != nil {
		return err
	}
	jwter, err := auth.NewJWTer(tokenCache, clocker)
	if err != nil {
		return err
	}
	// トランザクション
	appConnection := repository.NewAppConnection(db)

	// ルーティング設定
	groupRoute := router.Group("/api/v1").Use(handler.AuthMiddleware(jwter))

	getUsersService := service.NewGetUsers(db, rep, jwter)
	getUsersHandler := handler.NewGetUsers(getUsersService)
	groupRoute.GET("/users", getUsersHandler.ServeHTTP)

	getAccountService := service.NewGetAccount(db, rep)
	getUserHandler := handler.NewGetAccount(getAccountService)
	groupRoute.GET("/account", getUserHandler.ServeHTTP)

	signoutService := service.NewSignout(tokenCache)
	signoutHandler := handler.NewSignoutHandler(signoutService)
	groupRoute.DELETE("/signout", signoutHandler.ServeHTTP)

	sendPointService := service.NewSendPoint(rep, appConnection, db)
	sendPointHandler := handler.NewSendPoint(sendPointService)
	groupRoute.POST("/point_transactions", sendPointHandler.ServeHTTP)

	updatePassService := service.NewUpdatePassword(db, rep)
	updatePassHandler := handler.NewUpdatePasswordHandler(updatePassService)
	groupRoute.PATCH("/password", updatePassHandler.ServeHTTP)

	updateAccountService := service.NewUpdateAccount(db, rep)
	updateAccountHandler := handler.NewUpdateAccountHandler(updateAccountService)
	groupRoute.PUT("/account", updateAccountHandler.ServeHTTP)

	updateTemporaryEmailService := service.NewUpdateTemporaryEmail(db, cache, rep)
	updateTemporaryEmailHandler := handler.NewUpdateTemporaryEmailHandler(updateTemporaryEmailService)
	groupRoute.POST("/temporary_email", updateTemporaryEmailHandler.ServeHTTP)

	updateEmailService := service.NewUpdateEmail(db, cache, rep)
	updateEmailHandler := handler.NewUpdateEmailHandler(updateEmailService)
	groupRoute.PATCH("/email", updateEmailHandler.ServeHTTP)

	return nil
}
