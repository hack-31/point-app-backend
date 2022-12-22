package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hack-31/point-app-backend/auth"
)

// 認証ミドルウェア
func AuthMiddleware(j *auth.JWTer) gin.HandlerFunc {
	return gin.HandlerFunc(func(ctx *gin.Context) {
		log.Print("AuthMiddleware")

		if err := j.FillContext(ctx); err != nil {
			log.Print(err.Error())
			ErrResponse(ctx, http.StatusUnauthorized, "認証エラー", "アクセストークンが無効です。再ログインしてください。")

			return
		}
		ctx.Next()
	})
}
