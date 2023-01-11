package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
func (s *Signout) ServeHTTP(ctx *gin.Context) {
	const errTitle = "サインアウトエラー"
	// サインアウト処理依頼
	if err := s.Service.Signout(ctx); err != nil {
		ErrResponse(ctx, http.StatusInternalServerError, errTitle, err.Error())
		return
	}

	APIResponse(ctx, http.StatusOK, "サインアウトに成功しました。", nil)
}
