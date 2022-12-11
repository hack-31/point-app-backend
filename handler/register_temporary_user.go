package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/hack-31/point-app-backend/repository"
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
	var input struct {
		FirstName      string `json:"firstName"`
		FirstNameKana  string `json:"firstNameKana"`
		FamilyName     string `json:"familyName"`
		FamilyNameKana string `json:"familyNameKana"`
		Password       string `json:"password"`
		Email          string `json:"email"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		APIResponse(ctx, err.Error(), http.StatusInternalServerError, http.MethodPost, nil)
		return
	}

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
		APIResponse(ctx, err.Error(), http.StatusBadRequest, http.MethodPost, nil)
		return
	}

	err = ru.Service.RegisterTemporaryUser(ctx, input.FirstName, input.FirstNameKana, input.FamilyName, input.FamilyNameKana, input.Email, input.Password)

	if err != nil {
		if errors.Is(err, repository.ErrAlreadyEntry) {
			APIResponse(ctx, "登録済みのメールアドレスの登録できません。", http.StatusConflict, http.MethodPost, nil)
			return
		}

		APIResponse(ctx, err.Error(), http.StatusInternalServerError, http.MethodPost, nil)
		return
	}

	APIResponse(ctx, "本登録メールを送信しました。", http.StatusOK, http.MethodPost, nil)
}