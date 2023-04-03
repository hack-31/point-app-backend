package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/hack-31/point-app-backend/repository"
)

type RegisterTemporaryEmail struct {
	Service RegisterTemporaryEmailService
}

func NewRegisterTemporaryEmailHandler(s RegisterTemporaryEmailService) *RegisterTemporaryEmail {
	return &RegisterTemporaryEmail{Service: s}
}

// メール仮登録ハンドラー
//
// @param ctx ginContext
func (ru *RegisterTemporaryEmail) ServeHTTP(ctx *gin.Context) {
	const errTitle = "メールアドレス仮登録エラー"

	var input struct {
		Email string `json:"email"`
	}

	// ユーザーから正しいパラメータでポストデータが送られていない
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ErrResponse(ctx, http.StatusBadRequest, "パラメータエラー", err.Error())
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
		ErrResponse(ctx, http.StatusBadRequest, "パラメータエラー", err.Error())
		return
	}

	// サービス層にメール仮登録処理を依頼
	sessionID, err := ru.Service.RegisterTemporaryEmail(ctx, input.Email)

	// エラーレスポンスを返す
	if err != nil {
		if errors.Is(err, repository.ErrAlreadyEntry) {
			ErrResponse(ctx, http.StatusConflict, errTitle, repository.ErrAlreadyEntry.Error())
			return
		}
		ErrResponse(ctx, http.StatusInternalServerError, "サーバーエラー", err.Error())
		return
	}

	// 成功時のレスポンスを返す
	rsp := struct {
		Email string `json:"temporaryEmailId"`
	}{Email: sessionID}
	APIResponse(ctx, http.StatusCreated, "指定のメールアドレスに確認コードを送信しました。", rsp)
}
