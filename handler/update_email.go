package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hack-31/point-app-backend/repository"
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
	const mailErrTitle = "メールアドレス本更新エラー"
	const paramErrTitle = "パラメータエラー"
	const serverErrTitle = "サーバーエラー"

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
	if err := ue.Service.UpdateEmail(ctx, input.TemporaryEmailID, input.ConfirmCode); err != nil {
		// 確認コードとトークンが無効
		if errors.Is(err, repository.ErrNotFoundSession) {
			ErrResponse(ctx, http.StatusUnauthorized, mailErrTitle, repository.ErrNotFoundSession.Error())
			return
		}
		// 登録済みのメールアドレス
		if errors.Is(err, repository.ErrAlreadyEntry) {
			ErrResponse(ctx, http.StatusConflict, mailErrTitle, repository.ErrAlreadyEntry.Error())
			return
		}
		ErrResponse(ctx, http.StatusInternalServerError, serverErrTitle, err.Error())
		return
	}

	// 成功レスポンス
	APIResponse(ctx, http.StatusCreated, "更新が完了しました。", nil)
}
