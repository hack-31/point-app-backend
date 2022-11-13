package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/hack-31/point-app-backend/entity"
	"github.com/hack-31/point-app-backend/repository"
)

type RegisterUser struct {
	Service RegisterUserService
}

func NewRegisterUserHandler(s RegisterUserService) *RegisterUser {
	return &RegisterUser{Service: s}
}

// ユーザ登録ハンドラー
//
// @param ctx ginContext
func (ru *RegisterUser) ServeHTTP(ctx *gin.Context) {
	var input struct {
		Name     string `json:"name"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		APIResponse(ctx, err.Error(), http.StatusInternalServerError, http.MethodPost, nil)
		return
	}

	err := validation.ValidateStruct(&input,
		validation.Field(
			&input.Email,
			validation.Required,
			validation.Length(1, 20),
			is.Email,
		),
		validation.Field(
			&input.Name,
			validation.Required,
			validation.Length(1, 20),
		),
		validation.Field(
			&input.Password,
			validation.Required,
			validation.Length(1, 80),
		))
	if err != nil {
		APIResponse(ctx, err.Error(), http.StatusBadRequest, http.MethodPost, nil)
		return
	}

	// TODO: ユーザからadmin, generalなどの文字列を受けとって東麓する
	generalRole := 1
	u, err := ru.Service.RegisterUser(ctx, input.Name, input.Password, input.Email, generalRole)

	if err != nil {
		if errors.Is(err, repository.ErrAlreadyEntry) {
			APIResponse(ctx, "登録済みのメールアドレスです。", http.StatusConflict, http.MethodPost, nil)
			return
		}
		APIResponse(ctx, err.Error(), http.StatusInternalServerError, http.MethodPost, nil)
		return
	}

	rsp := struct {
		ID entity.UserID `json:"id"`
	}{ID: u.ID}
	APIResponse(ctx, "ユーザ登録に成功しました。", http.StatusOK, http.MethodPost, rsp)
}
