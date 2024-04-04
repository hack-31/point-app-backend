package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type UpdateAccount struct {
	Service UpdateAccountService
}

func NewUpdateAccountHandler(s UpdateAccountService) *UpdateAccount {
	return &UpdateAccount{Service: s}
}

// アカウント情報更新ハンドラー
//
// @param ctx ginContext
func (ua *UpdateAccount) ServeHTTP(ctx *gin.Context) {
	const errTitle = "パラメータエラー"

	// JSONを構造体にマッピング
	var input struct {
		FamilyName     string `json:"familyName"`
		FamilyNameKana string `json:"familyNameKana"`
		FirstName      string `json:"firstName"`
		FirstNameKana  string `json:"firstNameKana"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ErrResponse(ctx, http.StatusBadRequest, errTitle, err.Error(), err)
		return
	}

	// バリデーション
	err := validation.ValidateStruct(&input,
		validation.Field(
			&input.FamilyName,
			validation.Required,
			validation.Length(1, 50),
		),
		validation.Field(
			&input.FamilyNameKana,
			validation.Required,
			validation.Length(1, 50),
		),
		validation.Field(
			&input.FirstName,
			validation.Required,
			validation.Length(1, 50),
		),
		validation.Field(
			&input.FirstNameKana,
			validation.Required,
			validation.Length(1, 50),
		),
	)
	if err != nil {
		ErrResponse(ctx, http.StatusBadRequest, errTitle, err.Error(), err)
		return
	}

	// 更新を依頼
	err = ua.Service.UpdateAccount(ctx, input.FamilyName, input.FamilyNameKana, input.FirstName, input.FirstNameKana)

	if err != nil {
		ErrResponse(ctx, http.StatusInternalServerError, "サーバーエラー", err.Error(), err)
		return
	}

	// 成功レスポンス
	APIResponse(ctx, http.StatusCreated, "アカウント情報の更新に成功しました。", nil)
}
