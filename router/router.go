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

// 認証がないルーティングの設定を行う
//
// @param
// ctx コンテキスト
// router ルーター
func SetRouting(ctx context.Context, db *sqlx.DB, router *gin.Engine, cfg *config.Config) error {
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

	// ルーティング設定
	healthCheckhandler := handler.NewHealthCheckHandler()
	router.GET("/healthcheck", healthCheckhandler.ServeHTTP)

	groupRoute := router.Group("/api/v1")
	registerUserService := service.NewRegisterUser(db, rep, cache, jwter)
	registerUserHandler := handler.NewRegisterUserHandler(registerUserService)
	groupRoute.POST("/users", registerUserHandler.ServeHTTP)

	registerTempUserService := service.NewRegisterTemporaryUser(db, rep, cache)
	registerTempUserHandler := handler.NewRegisterTemporaryUserHandler(registerTempUserService)
	groupRoute.POST("/temporary_users", registerTempUserHandler.ServeHTTP)

	signinService := service.NewSignin(db, rep, cache, jwter)
	signinHandler := handler.NewSigninHandler(signinService)
	groupRoute.POST("/signin", signinHandler.ServeHTTP)

	resetPassService := service.NewResetPassword(db, rep)
	resetPassHandler := handler.NewResetPasswordHandler(resetPassService)
	groupRoute.PATCH("/random_password", resetPassHandler.ServeHTTP)

	return nil
}
