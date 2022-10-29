package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
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
	ctx.ShouldBindJSON(&input)
	
	// TODO: ユーザからadmin, generalなどの文字列を受けとって東麓する
	generalRole := 0
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
