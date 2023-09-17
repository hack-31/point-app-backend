package service

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hack-31/point-app-backend/auth"
	"github.com/hack-31/point-app-backend/domain/model"
	"github.com/hack-31/point-app-backend/repository"
	"github.com/hack-31/point-app-backend/utils/clock"
	"github.com/stretchr/testify/assert"
)

func TestGetNotifications(t *testing.T) {
	type input struct {
		nextToken string
		size      string
	}
	type want struct {
		ns  GetNotificationsResponse
		err error
	}
	type getNotifications struct {
		startID       model.NotificationID
		size          int
		notifications model.Notifications
		err           error
	}
	type getUserByID struct {
		user model.User
		err  error
	}

	userID := model.UserID(1)
	tests := map[string]struct {
		input            input
		getNotifications getNotifications
		getUserByID      getUserByID
		want             want
	}{
		"GetNotificationsサービスは、指定されたお知らせIDを開始位置として、サイズ数に応じたお知らせの一覧を返す": {
			input: input{
				nextToken: "90",
				size:      "5",
			},
			getNotifications: getNotifications{
				startID: model.NotificationID(90),
				size:    5,
				err:     nil,
				notifications: model.Notifications{
					&model.Notification{
						ID:        90,
						CreatedAt: clock.FixedClocker{}.Now(),
					},
					&model.Notification{
						ID:        81,
						CreatedAt: clock.FixedClocker{}.Now(),
					},
				},
			},
			getUserByID: getUserByID{
				err:  nil,
				user: model.User{},
			},
			want: want{
				ns: GetNotificationsResponse{
					NextToken: "80",
					Notifications: []struct {
						ID          model.NotificationID
						Title       string
						Description string
						IsChecked   bool
						CreatedAt   string
					}{
						{ID: 90, CreatedAt: "2022/05/10 12:34:56", Title: "", Description: "", IsChecked: false},
						{ID: 81, CreatedAt: "2022/05/10 12:34:56", Title: "", Description: "", IsChecked: false},
					},
				},
				err: nil,
			},
		},
		"お知らせ数が0件の時は、nextTokenは0、お知らせの配列は空配列を返す": {
			input: input{
				nextToken: "90",
				size:      "5",
			},
			getNotifications: getNotifications{
				startID:       model.NotificationID(90),
				size:          5,
				err:           nil,
				notifications: model.Notifications{},
			},
			getUserByID: getUserByID{
				err:  nil,
				user: model.User{},
			},
			want: want{
				ns: GetNotificationsResponse{
					NextToken: "0",
					Notifications: []struct {
						ID          model.NotificationID
						Title       string
						Description string
						IsChecked   bool
						CreatedAt   string
					}{},
				},
				err: nil,
			},
		},
		"nextTokenが空文字の場合は、ユーザテーブルより最新のお知らせIDを取得し、お知らせ一覧を返す": {
			input: input{
				nextToken: "",
				size:      "5",
			},
			getNotifications: getNotifications{
				startID: model.NotificationID(90),
				size:    5,
				err:     nil,
				notifications: model.Notifications{
					&model.Notification{
						ID:        90,
						CreatedAt: clock.FixedClocker{}.Now(),
					},
					&model.Notification{
						ID:        81,
						CreatedAt: clock.FixedClocker{}.Now(),
					},
				},
			},
			getUserByID: getUserByID{
				err: nil,
				user: model.User{
					ID:                   userID,
					NotificationLatestID: 90,
					Email:                "yamada@sample.com",
					SendingPoint:         100,
				},
			},
			want: want{
				ns: GetNotificationsResponse{
					NextToken: "80",
					Notifications: []struct {
						ID          model.NotificationID
						Title       string
						Description string
						IsChecked   bool
						CreatedAt   string
					}{
						{ID: 90, CreatedAt: "2022/05/10 12:34:56", Title: "", Description: "", IsChecked: false},
						{ID: 81, CreatedAt: "2022/05/10 12:34:56", Title: "", Description: "", IsChecked: false},
					},
				},
				err: nil,
			},
		},
	}

	for n, tt := range tests {
		tt := tt
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Set(auth.UserID, model.UserID(userID))

			// モックの定義
			moqQueryer := &QueryerMock{}
			moqNotificationRepo := &NotificationRepoMock{
				GetNotificationsFunc: func(ctx context.Context, db repository.Queryer, uid model.UserID, startID model.NotificationID, size int) (model.Notifications, error) {
					assert.Equal(t, tt.getNotifications.startID, startID)
					assert.Equal(t, tt.getNotifications.size, size)
					return tt.getNotifications.notifications, tt.getNotifications.err
				},
			}
			moqUserRepo := &UserRepoMock{
				GetUserByIDFunc: func(ctx context.Context, db repository.Queryer, ID model.UserID) (model.User, error) {
					assert.Equal(t, userID, ID)
					return tt.getUserByID.user, tt.getUserByID.err
				},
			}

			// サービス実行
			gns := &GetNotifications{
				DB:        moqQueryer,
				NotifRepo: moqNotificationRepo,
				UserRepo:  moqUserRepo,
			}
			gotNs, gotErr := gns.GetNotifications(ctx, tt.input.nextToken, tt.input.size)

			// アサーション
			assert.ErrorIs(t, gotErr, tt.want.err)
			assert.Equal(t, tt.want.ns, gotNs)
		})
	}
}
