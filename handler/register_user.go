package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/hack-31/point-app-backend/entity"
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
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		APIResponse(ctx, err.Error(), http.StatusBadRequest, http.MethodPost, nil)
		return
	}

	err := validation.ValidateStruct(&input,
		validation.Field(
			&input.TemporaryUserId,
			validation.Required,
		),
	)
	if err != nil {
		APIResponse(ctx, err.Error(), http.StatusBadRequest, http.MethodPost, nil)
		return
	}

	u, err := ru.Service.RegisterUser(ctx, input.TemporaryUserId)
	if err != nil {
		if errors.Is(err, repository.ErrNotFoundSession) {
			APIResponse(ctx, repository.ErrNotFoundSession.Error(), http.StatusUnauthorized, http.MethodPost, nil)
			return
		}
		if errors.Is(err, repository.ErrAlreadyEntry) {
			APIResponse(ctx, repository.ErrAlreadyEntry.Error(), http.StatusConflict, http.MethodPost, nil)
			return
		}
		APIResponse(ctx, err.Error(), http.StatusInternalServerError, http.MethodPost, nil)
		return
	}

	rsp := struct {
		ID entity.UserID `json:"userId"`
	}{ID: u.ID}
	APIResponse(ctx, "本登録が完了しました。", http.StatusCreated, http.MethodPost, rsp)
}
