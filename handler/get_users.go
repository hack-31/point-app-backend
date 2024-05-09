package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hack-31/point-app-backend/repository/entity"
)

type GetUsers struct {
	Service GetUsersService
}

func NewGetUsers(s GetUsersService) *GetUsers {
	return &GetUsers{Service: s}
}

// ユーザ一覧取得ハンドラー
//
// @param ctx ginContext
func (gu *GetUsers) ServeHTTP(ctx *gin.Context) {
	const errTitle = "ユーザ一覧取得エラー"

	// サービスにユーザ仮登録処理を依頼
	users, err := gu.Service.GetUsers(ctx)

	// エラーレスポンスを返す
	if err != nil {
		ErrResponse(ctx, http.StatusInternalServerError, errTitle, err.Error(), err)
		return
	}

	// レスポンスの作成
	type user struct {
		AcquisitionPoint int           `json:"acquisitionPoint"`
		Email            string        `json:"email"`
		FirstName        string        `json:"firstName"`
		FirstNameKana    string        `json:"firstNameKana"`
		FamilyName       string        `json:"familyName"`
		FamilyNameKana   string        `json:"familyNameKana"`
		ID               entity.UserID `json:"id"`
	}

	usersRes := make([]user, 0, len(users.Users))
	for _, v := range users.Users {
		usersRes = append(usersRes, user{
			ID:               v.ID,
			AcquisitionPoint: v.AcquisitionPoint,
			Email:            v.Email,
			FirstName:        v.FirstName,
			FirstNameKana:    v.FirstNameKana,
			FamilyName:       v.FamilyName,
			FamilyNameKana:   v.FamilyNameKana,
		})
	}

	rsp := struct {
		Users []user `json:"users"`
	}{
		Users: usersRes,
	}

	APIResponse(ctx, http.StatusOK, "取得成功しました。", rsp)
}
