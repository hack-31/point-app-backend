package service

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hack-31/point-app-backend/auth"
	"github.com/hack-31/point-app-backend/domain/model"
	"github.com/hack-31/point-app-backend/repository"
	"github.com/hack-31/point-app-backend/repository/entities"
	"github.com/hack-31/point-app-backend/repository/entity"
	"github.com/hack-31/point-app-backend/utils/clock"
	"github.com/stretchr/testify/assert"
)

func TestGetNotifications(t *testing.T) {
	t.Parallel()
	type input struct {
		nextToken string
		size      string
	}
	type want struct {
		ns  GetNotificationsResponse
		err error
	}
	type getByToUserByStartIdOrderByLatest struct {
		startID       entity.NotificationID
		size          int
		userID        model.UserID
		notifications entity.Notifications
		err           error
	}
	type getByToUserOrderByLatest struct {
		size          int
		userID        model.UserID
		notifications entity.Notifications
		err           error
	}

	type getUserByID struct {
		user entities.User
		err  error
	}

	userID := model.UserID(1)
	tests := map[string]struct {
		input                             input
		getByToUserByStartIdOrderByLatest getByToUserByStartIdOrderByLatest
		getByToUserOrderByLatest          getByToUserOrderByLatest
		getUserByID                       getUserByID
		want                              want
	}{
		"GetNotificationsサービスは、指定されたお知らせIDを開始位置として、サイズ数に応じたお知らせの一覧を返す": {
			input: input{
				nextToken: "90",
				size:      "5",
			},
			getByToUserByStartIdOrderByLatest: getByToUserByStartIdOrderByLatest{
				startID: entity.NotificationID(90),
				size:    5,
				userID:  userID,
				err:     nil,
				notifications: entity.Notifications{
					&entity.Notification{
						ID:        90,
						CreatedAt: clock.FixedClocker{}.Now(),
					},
					&entity.Notification{
						ID:        81,
						CreatedAt: clock.FixedClocker{}.Now(),
					},
				},
			},
			getByToUserOrderByLatest: getByToUserOrderByLatest{},
			getUserByID: getUserByID{
				err:  nil,
				user: entities.User{},
			},
			want: want{
				ns: GetNotificationsResponse{
					NextToken: "80",
					Notifications: []struct {
						ID          entity.NotificationID
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
			getByToUserByStartIdOrderByLatest: getByToUserByStartIdOrderByLatest{
				startID:       entity.NotificationID(90),
				size:          5,
				userID:        userID,
				err:           nil,
				notifications: entity.Notifications{},
			},
			getByToUserOrderByLatest: getByToUserOrderByLatest{},
			getUserByID: getUserByID{
				err:  nil,
				user: entities.User{},
			},
			want: want{
				ns: GetNotificationsResponse{
					NextToken: "0",
					Notifications: []struct {
						ID          entity.NotificationID
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
			getByToUserOrderByLatest: getByToUserOrderByLatest{
				userID: userID,
				size:   5,
				err:    nil,
				notifications: entity.Notifications{
					&entity.Notification{
						ID:        90,
						CreatedAt: clock.FixedClocker{}.Now(),
					},
					&entity.Notification{
						ID:        81,
						CreatedAt: clock.FixedClocker{}.Now(),
					},
				},
			},
			getUserByID: getUserByID{
				err: nil,
				user: entities.User{
					ID:           int64(userID),
					Email:        "yamada@sample.com",
					SendingPoint: 100,
				},
			},
			want: want{
				ns: GetNotificationsResponse{
					NextToken: "80",
					Notifications: []struct {
						ID          entity.NotificationID
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
				GetByToUserByStartIdOrderByLatestFunc: func(ctx context.Context, db repository.Queryer, uid model.UserID, startID entity.NotificationID, size int, columns ...string) (entity.Notifications, error) {
					assert.Equal(t, tt.getByToUserByStartIdOrderByLatest.startID, startID)
					assert.Equal(t, tt.getByToUserByStartIdOrderByLatest.size, size)
					assert.Equal(t, tt.getByToUserByStartIdOrderByLatest.userID, uid)
					return tt.getByToUserByStartIdOrderByLatest.notifications, tt.getByToUserByStartIdOrderByLatest.err
				},
				GetByToUserOrderByLatestFunc: func(ctx context.Context, db repository.Queryer, uid model.UserID, size int, columns ...string) (entity.Notifications, error) {
					assert.Equal(t, tt.getByToUserOrderByLatest.userID, uid)
					return tt.getByToUserOrderByLatest.notifications, tt.getByToUserByStartIdOrderByLatest.err
				},
			}
			moqUserRepo := &UserRepoMock{
				GetUserByIDFunc: func(ctx context.Context, db repository.Queryer, ID model.UserID) (entities.User, error) {
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
