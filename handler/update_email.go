package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type UpdateEmail struct {
	Service UpdateEmailService
}

func NewUpdateEmailHandler(s UpdateEmailService) *UpdateEmail {
	return &UpdateEmail{Service: s}
}

// メール本登録ハンドラー
//
// @param ctx ginContext
// @return APIレスポンス結果
func (ue *UpdateEmail) ServeHTTP(ctx *gin.Context) {
	// const mailErrTitle = "メールアドレス更新エラー"
	const paramErrTitle = "パラメータエラー"

	var input struct {
		TemporaryEmailID string `json:"temporaryEmailId"`
		ConfirmCode      string `json:"confirmCode"`
	}

	// ユーザーから正しいパラメータで送られているか確認
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ErrResponse(ctx, http.StatusBadRequest, paramErrTitle, err.Error())
		return
	}

	// バリデーションチェック
	err := validation.ValidateStruct(&input,
		validation.Field(
			&input.TemporaryEmailID,
			validation.Required,
		),
		validation.Field(
			&input.ConfirmCode,
			validation.Required,
			validation.Length(4, 4),
		),
	)
	if err != nil {
		ErrResponse(ctx, http.StatusBadRequest, paramErrTitle, err.Error())
		return
	}

	// サービス層に依頼する
	serviceErr := ue.Service.UpdateEmail(ctx, input.TemporaryEmailID, input.ConfirmCode)
	if err != nil {
		// TODO: サービス層エラーハンドリング
		ErrResponse(ctx, http.StatusBadRequest, paramErrTitle, serviceErr.Error())
		return
	}

	// 成功レスポンス
	APIResponse(ctx, http.StatusCreated, "更新が完了しました。", nil)

}
