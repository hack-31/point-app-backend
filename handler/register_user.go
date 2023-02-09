package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hack-31/point-app-backend/domain/model"
	"github.com/hack-31/point-app-backend/repository"
)

type RegisterUser struct {
	Service RegisterUserService
}

func NewRegisterUserHandler(s RegisterUserService) *RegisterUser {
	return &RegisterUser{Service: s}
}

// ユーザ本登録ハンドラー
//
// @param ctx ginContext
func (ru *RegisterUser) ServeHTTP(ctx *gin.Context) {
	var input struct {
		TemporaryUserId string `json:"temporaryUserId"`
		ConfirmCode     string `json:"confirmCode"`
	}
	const errTitle = "ユーザ登録エラー"
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ErrResponse(ctx, http.StatusBadRequest, errTitle, err.Error())
		return
	}

	err := validation.ValidateStruct(&input,
		validation.Field(
			&input.TemporaryUserId,
			validation.Required,
		),
		validation.Field(
			&input.ConfirmCode,
			validation.Required,
			validation.Length(4, 4),
		),
	)
	if err != nil {
		ErrResponse(ctx, http.StatusBadRequest, errTitle, err.Error())
		return
	}

	u, jwt, err := ru.Service.RegisterUser(ctx, input.TemporaryUserId, input.ConfirmCode)
	if err != nil {
		if errors.Is(err, repository.ErrNotFoundSession) {
			ErrResponse(ctx, http.StatusUnauthorized, errTitle, repository.ErrNotFoundSession.Error())
			return
		}
		if errors.Is(err, repository.ErrAlreadyEntry) {
			ErrResponse(ctx, http.StatusConflict, errTitle, repository.ErrAlreadyEntry.Error())
			return
		}
		ErrResponse(ctx, http.StatusInternalServerError, errTitle, err.Error())
		return
	}

	rsp := struct {
		ID    model.UserID `json:"userId"`
		Token string       `json:"accessToken"`
	}{ID: u.ID, Token: jwt}
	APIResponse(ctx, http.StatusCreated, "本登録が完了しました。", rsp)
}
