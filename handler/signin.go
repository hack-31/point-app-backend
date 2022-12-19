package handler

import (
	"errors"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/hack-31/point-app-backend/repository"
)

type Signin struct {
	Service SigninService
}

func NewSigninHandler(s SigninService) *Signin {
	return &Signin{Service: s}
}

// サインインハンドラー
//
// @param ctx ginContext
func (ru *Signin) ServeHTTP(ctx *gin.Context) {
	// ユーザのパラメータ検証
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		APIResponse(ctx, err.Error(), http.StatusBadRequest, http.MethodPost, nil)
		return
	}
	err := validation.ValidateStruct(&input,
		validation.Field(
			&input.Email,
			validation.Length(1, 256),
			validation.Required,
			is.Email,
		),
		validation.Field(
			&input.Password,
			validation.Length(8, 50),
			// TODO: 正規表現バリデーション
			validation.Match(regexp.MustCompile(``)),
			validation.Required,
		),
	)
	if err != nil {
		APIResponse(ctx, err.Error(), http.StatusBadRequest, http.MethodPost, nil)
		return
	}

	// サインイン処理依頼
	jwt, err := ru.Service.Signin(ctx, input.Email, input.Password)
	if err != nil {
		if errors.Is(err, repository.ErrNotMatchLogInfo) {
			APIResponse(ctx, repository.ErrNotMatchLogInfo.Error(), http.StatusUnauthorized, http.MethodPost, nil)
			return
		}
		APIResponse(ctx, err.Error(), http.StatusInternalServerError, http.MethodPost, nil)
		return
	}

	// 成功レスポンスを返す
	rsp := struct {
		Token string `json:"accessToken"`
	}{Token: jwt}
	APIResponse(ctx, "サインイン成功しました。", http.StatusCreated, http.MethodPost, rsp)
}
