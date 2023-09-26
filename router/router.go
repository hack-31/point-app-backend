package router

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/hack-31/point-app-backend/auth"
	"github.com/hack-31/point-app-backend/config"
	"github.com/hack-31/point-app-backend/repository"
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
	groupRoute := router.Group("/api/v1")
	router.GET("/healthcheck", InitHealthCheck().ServeHTTP)
	groupRoute.POST("/users", InitRegisterUser(db, rep, cache, jwter).ServeHTTP)
	groupRoute.POST("/temporary_users", InitRegisterTemporaryUser(db, rep, cache).ServeHTTP)
	groupRoute.POST("/signin", InitSignin(db, rep, cache, jwter).ServeHTTP)
	groupRoute.PATCH("/random_password", InitResetPassword(db, rep).ServeHTTP)

	return nil
}
