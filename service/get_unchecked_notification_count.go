package service

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/hack-31/point-app-backend/domain"
	"github.com/hack-31/point-app-backend/repository"
	"github.com/hack-31/point-app-backend/utils"
	"github.com/jmoiron/sqlx"
)

type GetUncheckedNotificationCount struct {
	DB        repository.Queryer
	Cache     domain.Cache
	NotifRepo domain.NotificationRepo
}

func NewGetUncheckedNotificationCount(db *sqlx.DB, cache domain.Cache, repo domain.NotificationRepo) *GetUncheckedNotificationCount {
	return &GetUncheckedNotificationCount{DB: db, Cache: cache, NotifRepo: repo}
}

// お知らせ数を確認し、かつ、お知らせ通知をサブスクライブ
//
// @params ctx コンテキスト
//
// @return
// ユーザ一覧
func (gunc *GetUncheckedNotificationCount) GetUncheckedNotificationCount(ctx *gin.Context) (<-chan int, error) {
	// ユーザID確認
	userID := utils.GetUserID(ctx)

	// お知らせ数の確認
	notificationCntChan := make(chan int)
	cnt, err := gunc.NotifRepo.GetUncheckedNotificationCount(ctx, gunc.DB, userID)
	if err != nil {
		return notificationCntChan, fmt.Errorf("faild to get unchecked notification count from db: %w", err)
	}

	// お知らせ通知をサブスクライブ
	go func() {
		defer close(notificationCntChan)

		// 初回のお知らせ通知
		notificationCntChan <- cnt

		// お知らせ通知をサブスク
		payloadChan, err := gunc.Cache.Subscribe(ctx, fmt.Sprintf("notification:%d", userID))
		if err != nil {
			return
		}
		for {
			select {
			case <-ctx.Request.Context().Done(): // リクエストのコンテキストがキャンセルされたらゴルーチンを終了
				return
			case <-ctx.Done():
				return
			case _, ok := <-payloadChan:
				// payloadChanがクローズされたらゴルーチンを終了
				if !ok {
					return
				}
				// お知らせの通知が来たら、お知らせテーブルよりお知らせ数取得
				cnt, err := gunc.NotifRepo.GetUncheckedNotificationCount(ctx, gunc.DB, userID)
				if err != nil {
					return
				}
				notificationCntChan <- cnt
			}
		}
	}()

	return notificationCntChan, nil
}
