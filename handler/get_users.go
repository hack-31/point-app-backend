package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hack-31/point-app-backend/repository"
)

type GetUsers struct {
	Service GetUsersService
}

func NewGetUsers(s GetUsersService) *GetUsers {
	return &GetUsers{Service: s}
}

type user struct {
	AcquisitionPoint int    `json:"acquisitionPoint"`
	Email            string `json:"email"`
	FirstName        string `json:"firstName"`
	FirstNameKana    string `json:"firstNameKana"`
	FamilyName       string `json:"familyName"`
	FamilyNameKana   string `json:"familyNameKana"`
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
		if errors.Is(err, repository.ErrAlreadyEntry) {
			ErrResponse(ctx, http.StatusConflict, errTitle, repository.ErrAlreadyEntry.Error())
			return
		}
		ErrResponse(ctx, http.StatusInternalServerError, errTitle, err.Error())
		return
	}

	usersResponse := []user{}

	for _, u := range users {
		usersResponse = append(usersResponse, user{
			AcquisitionPoint: u.AcquisitionPoint,
			Email:            u.Email,
			FirstName:        u.FirstName,
			FirstNameKana:    u.FirstNameKana,
			FamilyName:       u.FamilyName,
			FamilyNameKana:   u.FamilyNameKana},
		)
	}

	rsp := struct {
		Users []user `json:"users"`
	}{Users: usersResponse}

	APIResponse(ctx, http.StatusCreated, "取得成功しました。", rsp)
}
