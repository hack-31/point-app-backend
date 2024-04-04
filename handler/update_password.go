package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hack-31/point-app-backend/repository"
)

type UpdatePassword struct {
	Service UpdatePasswordService
}

func NewUpdatePasswordHandler(s UpdatePasswordService) *UpdatePassword {
	return &UpdatePassword{Service: s}
}

// パスワード更新ハンドラー
//
// @param ctx ginContext
func (rt *UpdatePassword) ServeHTTP(ctx *gin.Context) {
	const errTitle = "パスワード変更エラー"
	// JSONを構造体にマッピング
	var input struct {
		OldPassword string `json:"oldPassword"`
		NewPassword string `json:"newpassword"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ErrResponse(ctx, http.StatusBadRequest, errTitle, err.Error(), err)
		return
	}

	// バリデーション
	err := validation.ValidateStruct(&input,
		validation.Field(
			&input.OldPassword,
			validation.Required,
			validation.Length(8, 50),
		),
		validation.Field(
			&input.NewPassword,
			validation.Required,
			validation.Length(8, 50),
		),
	)
	if err != nil {
		ErrResponse(ctx, http.StatusBadRequest, errTitle, err.Error(), err)
		return
	}
	if input.NewPassword == input.OldPassword {
		ErrResponse(ctx, http.StatusBadRequest, errTitle, "古いパスワードと新しいパスワードが同じです。", err)
		return
	}

	// 更新を依頼
	err = rt.Service.UpdatePassword(ctx, input.OldPassword, input.NewPassword)
	if err != nil {
		if errors.Is(err, repository.ErrDifferentPassword) {
			ErrResponse(ctx, http.StatusBadRequest, errTitle, repository.ErrDifferentPassword.Error(), err)
			return
		}
		ErrResponse(ctx, http.StatusInternalServerError, errTitle, err.Error(), err)
		return
	}

	// 成功レスポンス
	APIResponse(ctx, http.StatusCreated, "パスワードを更新しました。", nil)
}
