package handler

import (
	"net/http"

	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/hack-31/point-app-backend/myerror"
)

type RegisterTemporaryUser struct {
	Service RegisterTemporaryUserService
}

func NewRegisterTemporaryUserHandler(s RegisterTemporaryUserService) *RegisterTemporaryUser {
	return &RegisterTemporaryUser{Service: s}
}

// ユーザ仮登録ハンドラー
//
// @param ctx ginContext
func (ru *RegisterTemporaryUser) ServeHTTP(ctx *gin.Context) {
	const errTitle = "ユーザ仮登録エラー"

	var input struct {
		FirstName      string `json:"firstName"`
		FirstNameKana  string `json:"firstNameKana"`
		FamilyName     string `json:"familyName"`
		FamilyNameKana string `json:"familyNameKana"`
		Password       string `json:"password"`
		Email          string `json:"email"`
	}
	// ユーザから正しいパラメータでポストデータが送られていない
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ErrResponse(ctx, http.StatusBadRequest, errTitle, err.Error(), err)
		return
	}

	// バリデーション検証
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
		validation.Field(
			&input.Email,
			validation.Required,
			validation.Length(1, 256),
			is.Email,
		),
		validation.Field(
			&input.Password,
			validation.Required,
			validation.Length(1, 50),
		))
	if err != nil {
		ErrResponse(ctx, http.StatusBadRequest, errTitle, err.Error(), err)
		return
	}
	// サービスにユーザ仮登録処理を依頼
	sessionID, err := ru.Service.RegisterTemporaryUser(ctx, input.FirstName, input.FirstNameKana, input.FamilyName, input.FamilyNameKana, input.Email, input.Password)

	// エラーレスポンスを返す
	if err != nil {
		if errors.Is(err, myerror.ErrAlreadyEntry) {
			ErrResponse(ctx, http.StatusConflict, errTitle, "登録済みのメールアドレスは登録できません。", err)
			return
		}
		ErrResponse(ctx, http.StatusInternalServerError, errTitle, "サーバーで予期せぬエラーが発生しました。時間をおいて再度お試しください。", err)
		return
	}

	// 成功時のレスポンスを返す
	rsp := struct {
		ID string `json:"temporaryUserId"`
	}{ID: sessionID}
	APIResponse(ctx, http.StatusCreated, "本登録メールを送信しました。", rsp)
}
