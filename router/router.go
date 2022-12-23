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
	rep := repository.Repository{Clocker: clocker}
	// キャッシュ
	cache, err := repository.NewKVS(ctx, cfg, repository.JWT)
	if err != nil {
		return err
	}
	// トークン
	tokenCache, err := repository.NewKVS(ctx, cfg, repository.TemporaryUserRegister)
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
	registerHandler := handler.NewRegisterUserHandler(&service.RegisterUser{DB: db, Cache: cache, TokenGenerator: jwter, Repo: &rep})
	groupRoute.POST("/users", registerHandler.ServeHTTP)

	registerTempUser := handler.NewRegisterTemporaryUserHandler(&service.RegisterTemporaryUser{DB: db, Cache: cache, Repo: &rep})
	groupRoute.POST("/temporary_users", registerTempUser.ServeHTTP)

	signin := handler.NewSigninHandler(&service.Signin{DB: db, Cache: cache, Repo: &rep, TokenGenerator: jwter})
	groupRoute.POST("/tokens", signin.ServeHTTP)

	signout := handler.NewSignoutHandler(&service.Signout{Cache: tokenCache})
	groupRoute.DELETE("/tokens/:userId", signout.ServeHTTP)

	resetPassword := handler.NewResetPasswordHandler(&service.ResetPassword{ExecerDB: db, QueryerDB: db, Repo: &rep})
	groupRoute.PATCH("/random_password", resetPassword.ServeHTTP)

	return nil
}
