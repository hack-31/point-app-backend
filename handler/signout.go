package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/hack-31/point-app-backend/repository"
)

type Signout struct {
	Service SignoutService
}

func NewSignoutHandler(s SignoutService) *Signout {
	return &Signout{Service: s}
}

// サインアウトハンドラー
//
// @param ctx ginContext
func (ru *Signout) ServeHTTP(ctx *gin.Context) {
	// ユーザのパラメータ検証
	uid := ctx.Param("userId")

	const errTitle = "サインアウトエラー"
	uidInt, err := strconv.Atoi(uid)
	if err != nil {
		ErrResponse(ctx, http.StatusBadRequest, errTitle, err.Error())
		return
	}

	if err := validation.Validate(uidInt,
		validation.Min(1),
		validation.Required,
	); err != nil {
		ErrResponse(ctx, http.StatusBadRequest, errTitle, err.Error())
		return
	}

	// サインアウト処理依頼
	if err := ru.Service.Signout(ctx, uid); err != nil {
		if errors.Is(err, repository.ErrNotFoundSession) {
			ErrResponse(ctx, http.StatusUnauthorized, errTitle, repository.ErrNotFoundSession.Error())
			return
		}
		ErrResponse(ctx, http.StatusInternalServerError, errTitle, err.Error())
		return
	}

	APIResponse(ctx, http.StatusCreated, "サインアウトに成功しました。", nil)
}
