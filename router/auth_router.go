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
	rep := repository.Repository{Clocker: clocker}
	// トークン
	tokenCache, err := repository.NewKVS(ctx, cfg, repository.TemporaryUserRegister)
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

	getUsersHandler := handler.NewGetUsers(&service.GetUsers{DB: db, Repo: &rep})
	groupRoute.GET("/users", getUsersHandler.ServeHTTP)

	getUserHandler := handler.NewGetAccount(&service.GetAccount{DB: db, Repo: &rep})
	groupRoute.GET("/account", getUserHandler.ServeHTTP)

	signout := handler.NewSignoutHandler(&service.Signout{Cache: tokenCache})
	groupRoute.DELETE("/signout", signout.ServeHTTP)

	sendPointHandler := handler.NewSendPoint(&service.SendPoint{PointRepo: &rep, UserRepo: &rep, Connection: appConnection, DB: db})
	groupRoute.POST("/point_transactions", sendPointHandler.ServeHTTP)

	return nil
}
