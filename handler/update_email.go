package handler

import (
	"net/http"

	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hack-31/point-app-backend/myerror"
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
		ErrResponse(ctx, http.StatusBadRequest, paramErrTitle, err.Error(), err)
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
		ErrResponse(ctx, http.StatusBadRequest, paramErrTitle, err.Error(), err)
		return
	}

	// サービス層に依頼する
	if err := ue.Service.UpdateEmail(ctx, input.TemporaryEmailID, input.ConfirmCode); err != nil {
		// 確認コードとトークンが無効
		if errors.Is(err, myerror.ErrNotFoundSession) {
			ErrResponse(ctx, http.StatusUnauthorized, mailErrTitle, myerror.ErrNotFoundSession.Error(), err)
			return
		}
		// 登録済みのメールアドレス
		if errors.Is(err, myerror.ErrAlreadyEntry) {
			ErrResponse(ctx, http.StatusConflict, mailErrTitle, myerror.ErrAlreadyEntry.Error(), err)
			return
		}
		ErrResponse(ctx, http.StatusInternalServerError, serverErrTitle, err.Error(), err)
		return
	}

	// 成功レスポンス
	APIResponse(ctx, http.StatusCreated, "更新が完了しました。", nil)
}
