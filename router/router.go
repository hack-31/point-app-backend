package router

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/hack-31/point-app-backend/handler"
	"github.com/hack-31/point-app-backend/repository"
	"github.com/hack-31/point-app-backend/service"
	"github.com/hack-31/point-app-backend/utils/clock"
	"github.com/jmoiron/sqlx"
)

// ルーティングの設定を行う
//
// @param ctx コンテキスト
//
// @param router ルーター
func SetRouting(ctx context.Context, db *sqlx.DB, router *gin.Engine) {

	clocker := clock.RealClocker{}
	rep := repository.Repository{Clocker: clocker}

	groupRoute := router.Group("/api/v1")

	registerHandler := handler.NewRegisterUserHandler(&service.RegisterUser{DB: db, Repo: &rep})
	groupRoute.POST("/register", registerHandler.ServeHTTP)
}
