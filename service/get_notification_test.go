package service

import (
	"database/sql"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hack-31/point-app-backend/auth"
	mock_domain "github.com/hack-31/point-app-backend/domain/_mock"
	"github.com/hack-31/point-app-backend/domain/model"
	"github.com/hack-31/point-app-backend/myerror"
	mock_repository "github.com/hack-31/point-app-backend/repository/_mock"
	customentities "github.com/hack-31/point-app-backend/repository/custom_entities"
	"github.com/hack-31/point-app-backend/repository/entities"
	"github.com/hack-31/point-app-backend/testutil"
	"github.com/hack-31/point-app-backend/utils/clock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
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
		callCount int
		err       error
	}
	type getNotificationID struct {
		callCount    int
		notification customentities.Notification
		err          error
	}
	type publish struct {
		callCount int
		err       error
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
				callCount: 1,
				err:       nil,
			},
			getNotificationByID: getNotificationID{
				callCount: 1,
				err:       nil,
				notification: customentities.Notification{
					Notification: entities.Notification{
						ID:          1,
						Description: "ポイント送付された",
						IsChecked:   false,
						CreatedAt:   clock.FixedClocker{}.Now(),
					},
					Type: entities.NotificationType{
						Title: "お知らせ",
					},
				},
			},
			publish: publish{
				callCount: 1,
				err:       nil,
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
				callCount: 1,
				err:       sql.ErrConnDone,
			},
			getNotificationByID: getNotificationID{
				callCount:    0,
				err:          nil,
				notification: customentities.Notification{},
			},
			publish: publish{
				callCount: 0,
				err:       nil,
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
				callCount: 1,
				err:       nil,
			},
			getNotificationByID: getNotificationID{
				callCount:    1,
				err:          sql.ErrConnDone,
				notification: customentities.Notification{},
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
				callCount: 1,
				err:       nil,
			},
			getNotificationByID: getNotificationID{
				callCount: 1,
				err:       nil,
				notification: customentities.Notification{
					Notification: entities.Notification{
						ID:          1,
						Description: "ポイント送付された",
						IsChecked:   false,
						CreatedAt:   clock.FixedClocker{}.Now(),
					},
					Type: entities.NotificationType{
						Title: "お知らせ",
					},
				},
			},
			publish: publish{
				callCount: 1,
				err:       myerror.ErrCacheException,
			},
			want: want{
				notification: GetNotificationResponse{
					ID:          1,
					Title:       "お知らせ",
					Description: "ポイント送付された",
					IsChecked:   false,
					CreatedAt:   "2022/05/10 12:34:56",
				},
				err: myerror.ErrCacheException,
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
			ctrl := gomock.NewController(t)
			mockTx := testutil.NewTxForMock(t, ctx)

			mockBeginner := mock_repository.NewMockBeginner(ctrl)
			mockBeginner.
				EXPECT().
				BeginTxx(ctx, nil).
				Return(mockTx, nil)

			mockCache := mock_domain.NewMockCache(ctrl)
			mockCache.
				EXPECT().
				Publish(ctx, fmt.Sprintf("notification:%d", tt.input.notificationID), gomock.Any()).
				Return(tt.publish.err).
				Times(tt.publish.callCount)

			mockNotifRepo := mock_domain.NewMockNotificationRepo(ctrl)
			mockNotifRepo.
				EXPECT().
				CheckNotification(ctx, mockTx, model.UserID(1), tt.input.notificationID).
				Return(tt.checkNotification.err).
				Times(tt.checkNotification.callCount)
			mockNotifRepo.
				EXPECT().
				GetNotificationByID(ctx, mockTx, model.UserID(1), tt.input.notificationID).
				Return(tt.getNotificationByID.notification, tt.getNotificationByID.err).
				Times(tt.getNotificationByID.callCount)

			// サービス実行
			gn := &GetNotification{
				Cache:     mockCache,
				Tx:        mockBeginner,
				NotifRepo: mockNotifRepo,
			}
			gotNs, gotErr := gn.GetNotification(ctx, tt.input.notificationID)

			// アサーション
			assert.ErrorIs(t, gotErr, tt.want.err)
			assert.Equal(t, tt.want.notification, gotNs)
		})
	}
}
