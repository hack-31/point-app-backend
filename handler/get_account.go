package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hack-31/point-app-backend/domain/model"
)

type GetAccount struct {
	Service GetAccountService
}

func NewGetAccount(s GetAccountService) *GetAccount {
	return &GetAccount{Service: s}
}

// ユーザ取得ハンドラー
//
// @param ctx ginContext
func (gu *GetAccount) ServeHTTP(ctx *gin.Context) {
	user, err := gu.Service.GetAccount(ctx)

	// エラーレスポンスを返す
	const errTitle = "アカウントエラー"
	if err != nil {
		ErrResponse(ctx, http.StatusInternalServerError, errTitle, err.Error())
		return
	}

	rsp := struct {
		AcquisitionPoint int          `json:"acquisitionPoint"`
		Email            string       `json:"email"`
		FamilyName       string       `json:"familyName"`
		FamilyNameKana   string       `json:"familyNameKana"`
		FirstName        string       `json:"firstName"`
		FirstNameKana    string       `json:"firstNameKana"`
		SendablePoint    int          `json:"sendablePoint"`
		UserID           model.UserID `json:"userId"`
	}{
		AcquisitionPoint: user.AcquisitionPoint,
		Email:            user.Email,
		FamilyName:       user.FamilyName,
		FamilyNameKana:   user.FamilyNameKana,
		FirstName:        user.FirstName,
		FirstNameKana:    user.FirstNameKana,
		UserID:           user.ID,
		SendablePoint:    user.SendablePoint,
	}

	APIResponse(ctx, http.StatusOK, "アカウント情報の取得に成功しました。", rsp)
}
