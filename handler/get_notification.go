package handler

import (
	"net/http"
	"strconv"

	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hack-31/point-app-backend/myerror"
	"github.com/hack-31/point-app-backend/repository/entity"
)

type GetNotification struct {
	Service GetNotificationService
}

func NewGetNotification(s GetNotificationService) *GetNotification {
	return &GetNotification{Service: s}
}

// お知らせ詳細取得ハンドラー
//
// @param ctx ginContext
func (gn *GetNotification) ServeHTTP(ctx *gin.Context) {
	const errTitle = "お知らせ取得エラー"

	// バリデーション検証
	ID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ErrResponse(ctx, http.StatusBadRequest, errTitle, "IDは数値を指定してください。", err)
		return
	}
	notificationID := entity.NotificationID(ID)
	if err := validation.Validate(notificationID, validation.Min(1), validation.Required); err != nil {
		ErrResponse(ctx, http.StatusBadRequest, errTitle, err.Error(), err)
		return
	}

	//　お知らせ詳細の取得
	n, err := gn.Service.GetNotification(ctx, entity.NotificationID(notificationID))
	if err != nil {
		if errors.Is(err, myerror.ErrNotFound) {
			ErrResponse(ctx, http.StatusNotFound, errTitle, myerror.ErrNotFound.Error(), err)
			return
		}
		ErrResponse(ctx, http.StatusInternalServerError, errTitle, err.Error(), err)
		return
	}

	// レスポンス作成
	rsp := struct {
		ID          entity.NotificationID `json:"id"`
		Title       string                `json:"title"`
		Description string                `json:"description"`
		IsChecked   bool                  `json:"isChecked"`
		CreatedAt   string                `json:"createdAt"`
	}{
		ID:          n.ID,
		Title:       n.Title,
		Description: n.Description,
		IsChecked:   n.IsChecked,
		CreatedAt:   n.CreatedAt,
	}
	APIResponse(ctx, http.StatusOK, "お知らせ情報の取得に成功しました。", rsp)
}
