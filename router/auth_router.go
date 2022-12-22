package router

import (
	"context"
	"net/http"

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
	// トークン
	clocker := clock.RealClocker{}
	tokenCache, err := repository.NewKVS(ctx, cfg, repository.TemporaryUserRegister)
	if err != nil {
		return err
	}
	jwter, err := auth.NewJWTer(tokenCache, clocker)
	if err != nil {
		return err
	}

	// ルーティング設定
	groupRoute := router.Group("/api/v1").Use(handler.AuthMiddleware(jwter))
	// TODO: 一旦仮の値を返す
	groupRoute.GET("/users", func(ctx *gin.Context) {
		email, _ := ctx.Get(auth.Email)
		type user struct {
			AcquisitionPoint int    `json:"acquisitionPoint"`
			Email            string `json:"email"`
			FirstName        string `json:"firstName"`
			FirstNameKana    string `json:"firstNameKana"`
			FamilyName       string `json:"familyName"`
			FamilyNameKana   string `json:"familyNameKana"`
		}
		u := []user{}
		u = append(u, user{
			AcquisitionPoint: 2000,
			Email:            email.(string),
			FirstName:        "田中",
			FirstNameKana:    "タナカ",
			FamilyName:       "和樹",
			FamilyNameKana:   "カズキ"},
		)
		rsp := struct {
			Users []user `json:"users"`
		}{Users: u}
		handler.APIResponse(ctx, http.StatusCreated, "認証成功", rsp)
	})
	return nil
}
