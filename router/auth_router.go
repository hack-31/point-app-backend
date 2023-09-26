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
	groupRoute := router.Group("/api/v1").
		Use(handler.AuthMiddleware(jwter))

	groupRoute.GET("/users", InitGetUsers(db, rep, jwter).ServeHTTP)
	groupRoute.GET("/account", InitGetAccount(db, rep).ServeHTTP)
	groupRoute.DELETE("/signout", InitSignout(tokenCache).ServeHTTP)
	groupRoute.POST("/point_transactions", InitSendPoint(rep, appConnection, cache).ServeHTTP)
	groupRoute.PATCH("/password", InitUpdatePassword(db, rep).ServeHTTP)
	groupRoute.PUT("/account", InitUpdateAccount(db, rep).ServeHTTP)
	groupRoute.POST("/temporary_email", InitUpdateTemporaryEmail(db, cache, rep).ServeHTTP)
	groupRoute.PATCH("/email", InitUpdateEmail(db, cache, rep).ServeHTTP)
	groupRoute.GET("/notifications/:id", InitGetNotification(cache, rep, appConnection).ServeHTTP)
	groupRoute.GET("/notifications", InitGetNotifications(db, rep).ServeHTTP)
	groupRoute.GET("/unchecked_notification_count", InitGetUncheckedNotificationCount(db, cache, rep).ServeHTTP)

	return nil
}
