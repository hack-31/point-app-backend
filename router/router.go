package router

import (
	"context"

	"github.com/gin-gonic/gin"
)

// ルーティングの設定を行う
//
// @param ctx コンテキスト
//
// @param router ルーター
func SetRouting(ctx context.Context, router *gin.Engine) {
	groupRoute := router.Group("/api/v1")
	// TODO: 一旦仮で固定値を返す
	groupRoute.GET("", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World",
		})
	})
}
