package handler

import (
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

type RegisterTemporaryEmail struct {
	// Service RegisterTemporaryEmailService
}

func NewRegisterTemporaryEmailHandler() *RegisterTemporaryEmail {
	return &RegisterTemporaryEmail{}
}

// メール仮登録ハンドラー
//
// @param ctx ginContext
func (ru *RegisterTemporaryEmail) ServeHTTP(ctx *gin.Context) {
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		println("エラー")
	}
	println(string(body))
}
