package service

import (
	"context"
	"database/sql"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hack-31/point-app-backend/auth"
	"github.com/hack-31/point-app-backend/domain/model"
	"github.com/hack-31/point-app-backend/repository"
	"github.com/hack-31/point-app-backend/utils/clock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestGetNotification(t *testing.T) {
	t.Parallel()
	type input struct {
		notificationID model.NotificationID
	}
	type want struct {
		notification GetNotificationResponse
		err          error
	}
	type checkNotification struct {
		err error
	}
	type getNotificationID struct {
		notification model.Notification
		err          error
	}
	type publish struct {
		err error
	}

	tests := map[string]struct {
		input               input
		checkNotification   checkNotification
		getNotificationByID getNotificationID
		want                want
		publish             publish
	}{
		"GetNotificationサービスはお知らせIDを渡すと、そのIDに応じたお知らせ詳細を返す": {
			input: input{
				notificationID: 1,
			},
			checkNotification: checkNotification{
				err: nil,
			},
			getNotificationByID: getNotificationID{
				err: nil,
				notification: model.Notification{
					ID:          1,
					Title:       "お知らせ",
					Description: "ポイント送付された",
					IsChecked:   false,
					CreatedAt:   clock.FixedClocker{}.Now(),
				},
			},
			want: want{
				notification: GetNotificationResponse{
					ID:          1,
					Title:       "お知らせ",
					Description: "ポイント送付された",
					IsChecked:   false,
					CreatedAt:   "2022/05/10 12:34:56",
				},
				err: nil,
			},
		},
		"リポジトリに対してお知らせを既読状態にする際、予期せぬエラーが発生した場合は、DB予期せぬエラーを返す": {
			input: input{
				notificationID: 1,
			},
			checkNotification: checkNotification{
				err: sql.ErrConnDone,
			},
			getNotificationByID: getNotificationID{
				err:          nil,
				notification: model.Notification{},
			},
			want: want{
				notification: GetNotificationResponse{},
				err:          sql.ErrConnDone,
			},
		},
		"リポジトリに対してお知らせ詳細を取得する際、予期せぬエラーが発生した場合は、DB予期せぬエラーを返す": {
			input: input{
				notificationID: 1,
			},
			checkNotification: checkNotification{
				err: nil,
			},
			getNotificationByID: getNotificationID{
				err:          sql.ErrConnDone,
				notification: model.Notification{},
			},
			want: want{
				notification: GetNotificationResponse{},
				err:          sql.ErrConnDone,
			},
		},
		"チャネルに対するpublishが失敗すると、キャッシュ予期せぬエラーとユーザ情報を返す": {
			input: input{
				notificationID: 1,
			},
			checkNotification: checkNotification{
				err: nil,
			},
			getNotificationByID: getNotificationID{
				err: nil,
				notification: model.Notification{
					ID:          1,
					Title:       "お知らせ",
					Description: "ポイント送付された",
					IsChecked:   false,
					CreatedAt:   clock.FixedClocker{}.Now(),
				},
			},
			publish: publish{
				err: repository.ErrCacheException,
			},
			want: want{
				notification: GetNotificationResponse{
					ID:          1,
					Title:       "お知らせ",
					Description: "ポイント送付された",
					IsChecked:   false,
					CreatedAt:   "2022/05/10 12:34:56",
				},
				err: repository.ErrCacheException,
			},
		},
	}

	for n, tt := range tests {
		tt := tt
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Set(auth.UserID, model.UserID(1))

			// モックの定義
			moqCache := &CacheMock{
				PublishFunc: func(ctx context.Context, channel, palyload string) error {
					wantChannel := fmt.Sprintf("notification:%d", tt.input.notificationID)
					assert.Equal(t, wantChannel, channel)
					return tt.publish.err
				},
			}
			moqTranactor := &TransacterMock{
				BeginFunc: func(ctx context.Context) error {
					return nil
				},
				CommitFunc: func() error {
					return nil
				},
				RollbackFunc: func() error {
					return nil
				},
				DBFunc: func() *sqlx.Tx {
					return &sqlx.Tx{}
				},
			}
			moqNotificationRepo := &NotificationRepoMock{
				CheckNotificationFunc: func(ctx context.Context, db repository.Execer, uid model.UserID, nid model.NotificationID) error {
					// TODO: 引数のテストも入れる
					return tt.checkNotification.err
				},
				GetNotificationByIDFunc: func(ctx context.Context, db repository.Queryer, uid model.UserID, nid model.NotificationID) (model.Notification, error) {
					// TODO: 引数のテストも入れる
					return tt.getNotificationByID.notification, tt.getNotificationByID.err
				},
			}

			// サービス実行
			gn := &GetNotification{
				Cache:      moqCache,
				Connection: moqTranactor,
				NotifRepo:  moqNotificationRepo,
			}
			gotNs, gotErr := gn.GetNotification(ctx, tt.input.notificationID)

			// アサーション
			assert.ErrorIs(t, gotErr, tt.want.err)
			assert.Equal(t, tt.want.notification, gotNs)
		})
	}
}
