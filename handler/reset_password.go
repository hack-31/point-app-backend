package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/hack-31/point-app-backend/repository"
)

type ResetPassword struct {
	Service ResetPasswordService
}

func NewResetPasswordHandler(s ResetPasswordService) *ResetPassword {
	return &ResetPassword{Service: s}
}

// サインアウトハンドラー
//
// @param ctx ginContext
func (s *ResetPassword) ServeHTTP(ctx *gin.Context) {
	// ユーザのパラメータ検証
	var input struct {
		Email string `json:"email"`
	}

	const errTitle = "パスワードリセットエラー"
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ErrResponse(ctx, http.StatusBadRequest, errTitle, err.Error())
		return
	}
	if err := validation.ValidateStruct(&input,
		validation.Field(
			&input.Email,
			validation.Length(1, 256),
			validation.Required,
			is.Email,
		),
	); err != nil {
		ErrResponse(ctx, http.StatusBadRequest, errTitle, err.Error())
		return
	}

	// パスワード再発行処理依頼
	if err := s.Service.ResetPassword(ctx, input.Email); err != nil {
		if errors.Is(err, repository.ErrNotExistEmail) {
			ErrResponse(ctx, http.StatusNotFound, errTitle, repository.ErrNotExistEmail.Error())
			return
		}
		ErrResponse(ctx, http.StatusInternalServerError, errTitle, err.Error())
		return
	}

	APIResponse(ctx, http.StatusCreated, "パスワード再発行メールを送信しました。", nil)
}
