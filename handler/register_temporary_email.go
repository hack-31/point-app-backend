package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
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
	const errTitle = "メール仮登録エラー"

	var input struct {
		Email string `json:"email"`
	}

	// ユーザーから正しいパラメータでポストデータが送られていない
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ErrResponse(ctx, http.StatusBadRequest, errTitle, err.Error())
		return
	}

	// バリデーション検証
	err := validation.ValidateStruct(&input,
		validation.Field(
			&input.Email,
			validation.Required,
			validation.Length(1, 256),
			is.Email,
		))
	if err != nil {
		ErrResponse(ctx, http.StatusBadRequest, errTitle, err.Error())
		return
	}
}
