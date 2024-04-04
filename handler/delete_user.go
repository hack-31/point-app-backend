package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hack-31/point-app-backend/domain/model"
	"github.com/hack-31/point-app-backend/repository"
)

type DeleteUser struct {
	Service DeleteUserService
}

func NewDeleteUser(s DeleteUserService) *DeleteUser {
	return &DeleteUser{Service: s}
}

// ユーザー削除取得ハンドラー
//
// @param ctx ginContext
func (du *DeleteUser) ServeHTTP(ctx *gin.Context) {
	const errTitle = "ユーザー削除エラー"

	// バリデーション検証
	ID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ErrResponse(ctx, http.StatusBadRequest, errTitle, "IDは数値を指定してください。", err)
		return
	}
	userID := model.UserID(ID)
	if err := validation.Validate(userID, validation.Min(1), validation.Required); err != nil {
		ErrResponse(ctx, http.StatusBadRequest, errTitle, err.Error(), err)
		return
	}

	if err := du.Service.DeleteUser(ctx, userID); err != nil {
		if errors.Is(err, repository.ErrNotUser) {
			ErrResponse(ctx, http.StatusNotFound, errTitle, repository.ErrNotFound.Error(), err)
			return
		}
		ErrResponse(ctx, http.StatusInternalServerError, errTitle, err.Error(), err)
		return
	}

	// レスポンス作成
	rsp := struct {
		ID model.UserID `json:"userId"`
	}{
		ID: userID,
	}
	APIResponse(ctx, http.StatusCreated, "ユーザー削除しました。", rsp)
}
