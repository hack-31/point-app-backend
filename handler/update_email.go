package handler

import "github.com/gin-gonic/gin"

type UpdateEmail struct{}

func NewUpdateEmailHandler() *UpdateEmail {
	return &UpdateEmail{}
}

// メール本登録ハンドラー
//
// @param ctx ginContext
func (ue *UpdateEmail) ServeHTTP(ctx *gin.Context) {
	println("テスト")
}
