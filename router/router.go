package router

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/hack-31/point-app-backend/auth"
	"github.com/hack-31/point-app-backend/config"
	"github.com/hack-31/point-app-backend/handler"
	"github.com/hack-31/point-app-backend/repository"
	"github.com/hack-31/point-app-backend/utils/clock"
	"github.com/jmoiron/sqlx"
)

// ルーティングの設定を行う
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
	// トランザクション
	transacter := repository.NewAppConnection(db)

	// ルーティング設定
	rootRoute := router.Group("/api/v1")
	publicRoute := rootRoute
	protectRoute := publicRoute.Use(handler.AuthMiddleware(jwter))

	router.GET("/healthcheck", InitHealthCheck().ServeHTTP)
	publicRoute.POST("/users", InitRegisterUser(db, rep, cache, jwter).ServeHTTP)
	publicRoute.POST("/temporary_users", InitRegisterTemporaryUser(db, rep, cache).ServeHTTP)
	publicRoute.POST("/signin", InitSignin(db, rep, cache, jwter).ServeHTTP)
	publicRoute.PATCH("/random_password", InitResetPassword(db, rep).ServeHTTP)

	protectRoute.GET("/users", InitGetUsers(db, rep, jwter).ServeHTTP)
	protectRoute.GET("/account", InitGetAccount(db, rep).ServeHTTP)
	protectRoute.DELETE("/signout", InitSignout(tokenCache).ServeHTTP)
	protectRoute.POST("/point_transactions", InitSendPoint(rep, transacter, cache).ServeHTTP)
	protectRoute.PATCH("/password", InitUpdatePassword(db, rep).ServeHTTP)
	protectRoute.PUT("/account", InitUpdateAccount(db, rep).ServeHTTP)
	protectRoute.POST("/temporary_email", InitUpdateTemporaryEmail(db, cache, rep).ServeHTTP)
	protectRoute.PATCH("/email", InitUpdateEmail(db, cache, rep).ServeHTTP)
	protectRoute.GET("/notifications/:id", InitGetNotification(cache, rep, transacter).ServeHTTP)
	protectRoute.GET("/notifications", InitGetNotifications(db, rep).ServeHTTP)
	protectRoute.GET("/unchecked_notification_count", InitGetUncheckedNotificationCount(db, cache, rep).ServeHTTP)

	return nil
}
