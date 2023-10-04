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
	r := router.Group("/api/v1")
	authMiddleware := handler.AuthMiddleware(jwter)

	// パブリック
	router.GET("/healthcheck", InitHealthCheck().ServeHTTP)
	r.POST("/users", InitRegisterUser(db, rep, cache, jwter).ServeHTTP)
	r.POST("/temporary_users", InitRegisterTemporaryUser(db, rep, cache).ServeHTTP)
	r.POST("/signin", InitSignin(db, rep, cache, jwter).ServeHTTP)
	r.PATCH("/random_password", InitResetPassword(db, rep).ServeHTTP)

	// プロテクト
	r.GET("/users", authMiddleware, InitGetUsers(db, rep, jwter).ServeHTTP)
	r.GET("/account", authMiddleware, InitGetAccount(db, rep).ServeHTTP)
	r.DELETE("/signout", authMiddleware, InitSignout(tokenCache).ServeHTTP)
	r.DELETE("/users/:id", authMiddleware, InitDeleteUser(transacter, rep, tokenCache).ServeHTTP)
	r.POST("/point_transactions", authMiddleware, InitSendPoint(rep, transacter, cache).ServeHTTP)
	r.PATCH("/password", authMiddleware, InitUpdatePassword(db, rep).ServeHTTP)
	r.PUT("/account", authMiddleware, InitUpdateAccount(db, rep).ServeHTTP)
	r.POST("/temporary_email", authMiddleware, InitUpdateTemporaryEmail(db, cache, rep).ServeHTTP)
	r.PATCH("/email", authMiddleware, InitUpdateEmail(db, cache, rep).ServeHTTP)
	r.GET("/notifications/:id", authMiddleware, InitGetNotification(cache, rep, transacter).ServeHTTP)
	r.GET("/notifications", authMiddleware, InitGetNotifications(db, rep).ServeHTTP)
	r.GET("/unchecked_notification_count", authMiddleware, InitGetUncheckedNotificationCount(db, cache, rep).ServeHTTP)

	return nil
}
