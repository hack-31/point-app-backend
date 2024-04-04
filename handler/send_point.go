package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hack-31/point-app-backend/repository"
)

type SendPoint struct {
	Service SendPointService
}

func NewSendPoint(s SendPointService) *SendPoint {
	return &SendPoint{Service: s}
}

// ポイント送付ハンドラー
//
// @param ctx ginContext
func (sp *SendPoint) ServeHTTP(ctx *gin.Context) {
	// ユーザのパラメータ検証
	var input struct {
		ToUserId  int `json:"toUserId"`
		SendPoint int `json:"sendPoint"`
	}

	const errTitle = "送付エラー"
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ErrResponse(ctx, http.StatusBadRequest, errTitle, err.Error(), err)
		return
	}

	// バリデーション
	err := validation.ValidateStruct(&input,
		validation.Field(
			&input.ToUserId,
			validation.Required,
		),
		validation.Field(
			&input.SendPoint,
			validation.Required,
		),
	)
	if err != nil {
		ErrResponse(ctx, http.StatusBadRequest, errTitle, err.Error(), err)
		return
	}

	// ポイント送付
	if err := sp.Service.SendPoint(ctx, input.ToUserId, input.SendPoint); err != nil {
		if errors.Is(err, repository.ErrNotUser) {
			ErrResponse(ctx, http.StatusNotFound, errTitle, repository.ErrNotUser.Error(), err)
			return
		}
		if errors.Is(err, repository.ErrHasNotSendablePoint) {
			ErrResponse(ctx, http.StatusBadRequest, errTitle, repository.ErrHasNotSendablePoint.Error(), err)
			return
		}
		ErrResponse(ctx, http.StatusInternalServerError, errTitle, err.Error(), err)
		return
	}

	APIResponse(ctx, http.StatusCreated, "ポイントの送信に成功しました。", nil)
}
