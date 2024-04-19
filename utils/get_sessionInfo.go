package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/hack-31/point-app-backend/auth"
	"github.com/hack-31/point-app-backend/repository/entity"
)

// コンテキストからメールアドレスを取得
func GetEmail(ctx *gin.Context) string {
	if email, ok := ctx.Get(auth.Email); ok {
		return email.(string)
	}
	return ""
}

// コンテキストからユーザーIDを取得
func GetUserID(ctx *gin.Context) entity.UserID {
	if userID, ok := ctx.Get(auth.UserID); ok {
		return userID.(entity.UserID)
	}
	return entity.UserID(0)
}
