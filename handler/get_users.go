package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/hack-31/point-app-backend/domain/model"
	"github.com/hack-31/point-app-backend/service"
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

	// クエリの取得
	queries := struct {
		Size      string `json:"size"`
		NextToken string `json:"nextToken"`
	}{
		Size:      ctx.Query("size"),
		NextToken: ctx.Query("nextToken"),
	}

	// バリデーション
	err := validation.ValidateStruct(&queries,
		validation.Field(
			&queries.Size,
			is.Int.Error("数値で入力してください"),
			validation.Min(1).Error("1以上の数値で入力してください"),
		),
		validation.Field(
			&queries.NextToken,
			is.Base64.Error("base64形式で入力してください"),
		),
	)
	if err != nil {
		ErrResponse(ctx, http.StatusBadRequest, errTitle, err.Error(), err)
		return
	}

	// サービスにユーザ仮登録処理を依頼
	sizeInt, _ := strconv.Atoi(queries.Size)
	users, err := gu.Service.GetUsers(ctx, service.GetUsersRequest{
		Size:       sizeInt,
		NextCursor: queries.NextToken,
	})

	// エラーレスポンスを返す
	if err != nil {
		ErrResponse(ctx, http.StatusInternalServerError, errTitle, err.Error(), err)
		return
	}

	// レスポンスの作成
	type user struct {
		AcquisitionPoint int          `json:"acquisitionPoint"`
		Email            string       `json:"email"`
		FirstName        string       `json:"firstName"`
		FirstNameKana    string       `json:"firstNameKana"`
		FamilyName       string       `json:"familyName"`
		FamilyNameKana   string       `json:"familyNameKana"`
		ID               model.UserID `json:"id"`
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
		Users     []user `json:"users"`
		NextToken string `json:"nextToken"`
	}{
		Users:     usersRes,
		NextToken: users.NextCursor,
	}

	APIResponse(ctx, http.StatusOK, "取得成功しました。", rsp)
}
