package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/hack-31/point-app-backend/domain/model"
)

type GetNotifications struct {
	Service GetNotificationsService
}

func NewGetNotifications(s GetNotificationsService) *GetNotifications {
	return &GetNotifications{Service: s}
}

// お知らせ一覧取得ハンドラー
//
// @param ctx ginContext
func (gn *GetNotifications) ServeHTTP(ctx *gin.Context) {
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
			is.Int,
			validation.Required,
		),
		validation.Field(
			&queries.NextToken,
			is.Int,
		),
	)
	const errTitle = "お知らせ一覧取得エラー"
	if err != nil {
		ErrResponse(ctx, http.StatusBadRequest, errTitle, err.Error(), err)
		return
	}

	// お知らせ一覧取得
	ns, err := gn.Service.GetNotifications(ctx, queries.NextToken, queries.Size)
	if err != nil {
		ErrResponse(ctx, http.StatusInternalServerError, errTitle, err.Error(), err)
		return
	}

	// レスポンス作成
	type notification struct {
		ID          model.NotificationID `json:"id"`
		Title       string               `json:"title"`
		Description string               `json:"description"`
		IsChecked   bool                 `json:"isChecked"`
	}
	notifications := []notification{}
	for _, n := range ns.Notifications {
		notifications = append(notifications, notification{
			ID:          n.ID,
			Title:       n.Title,
			Description: n.Description,
			IsChecked:   n.IsChecked,
		})
	}
	nextToken, _ := strconv.Atoi(ns.NextToken)
	rsp := struct {
		Notifications []notification `json:"notifications"`
		NextToken     int            `json:"nextToken"`
	}{
		Notifications: notifications,
		NextToken:     nextToken,
	}
	APIResponse(ctx, http.StatusOK, "取得に成功しました。", rsp)
}
