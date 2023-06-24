package service

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/hack-31/point-app-backend/auth"
	"github.com/hack-31/point-app-backend/domain"
	"github.com/hack-31/point-app-backend/domain/model"
	"github.com/hack-31/point-app-backend/repository"
)

type GetNotification struct {
	Connection repository.Transacter
	Cache      domain.Cache
	NotifRepo  domain.NotificationRepo
}

func NewGetNotification(cache domain.Cache, repo *repository.Repository, connection repository.Transacter) *GetNotification {
	return &GetNotification{Cache: cache, NotifRepo: repo, Connection: connection}
}

type GetNotificationResponse struct {
	ID          model.NotificationID
	Title       string
	Description string
	IsChecked   bool
	CreatedAt   string
}

// お知らせ詳細取得サービス
//
// @params ctx コンテキスト
// @params notificationID お知らせID
//
// @return
// お知らせ詳細
func (gn *GetNotification) GetNotification(ctx *gin.Context, notificationID model.NotificationID) (GetNotificationResponse, error) {
	// ユーザID確認
	u, _ := ctx.Get(auth.UserID)
	userID := u.(model.UserID)

	// トランザクション開始
	if err := gn.Connection.Begin(ctx); err != nil {
		return GetNotificationResponse{}, err
	}
	// ロールバック
	defer func() { _ = gn.Connection.Rollback() }()

	// 閲覧したので、お知らせをチェック済みとする
	if err := gn.NotifRepo.CheckNotification(ctx, gn.Connection.DB(), userID, notificationID); err != nil {
		return GetNotificationResponse{}, fmt.Errorf("cannot check notification in db: %w", err)
	}

	// お知らせ詳細取得
	n, err := gn.NotifRepo.GetNotificationByID(ctx, gn.Connection.DB(), userID, notificationID)
	if err != nil {
		return GetNotificationResponse{}, fmt.Errorf("cannot GetNotificaitonByID in db: %w", err)
	}
	res := GetNotificationResponse{
		ID:          n.ID,
		Title:       n.Title,
		IsChecked:   n.IsChecked,
		Description: n.Description,
		CreatedAt:   model.NewTime(n.CreatedAt).Format(),
	}

	// トランザクションコミット
	if err := gn.Connection.Commit(); err != nil {
		return GetNotificationResponse{}, err
	}

	// お知らせチェックしたので、お知らせを通知
	channel := fmt.Sprintf("notification:%d", userID)
	payload, err := json.Marshal(n)
	if err != nil {
		return res, fmt.Errorf("cannot marshal: %w ", err)
	}
	if err := gn.Cache.Publish(ctx, channel, string(payload)); err != nil {
		return res, fmt.Errorf("cannot publish to %s channel: %w ", channel, err)
	}
	return res, nil
}
