package service

import (
	"strconv"

	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
	"github.com/hack-31/point-app-backend/domain"
	"github.com/hack-31/point-app-backend/domain/model"
	"github.com/hack-31/point-app-backend/repository"
	customentities "github.com/hack-31/point-app-backend/repository/custom_entities"
	"github.com/hack-31/point-app-backend/utils"
	"github.com/jmoiron/sqlx"
)

type GetNotifications struct {
	DB        repository.Queryer
	NotifRepo domain.NotificationRepo
	UserRepo  domain.UserRepo
}

func NewGetNotifications(db *sqlx.DB, repo *repository.Repository) *GetNotifications {
	return &GetNotifications{DB: db, NotifRepo: repo, UserRepo: repo}
}

type GetNotificationsResponse struct {
	Notifications []struct {
		ID          model.NotificationID
		Title       string
		Description string
		IsChecked   bool
		CreatedAt   string
	}
	NextToken string
}

// お知らせ一覧取得サービス
//
// @params ctx コンテキスト
//
// @return
// ユーザ一覧
func (gn *GetNotifications) GetNotifications(ctx *gin.Context, nextToken, size string) (GetNotificationsResponse, error) {
	// ユーザID確認
	userID := utils.GetUserID(ctx)

	var ns []*customentities.Notification
	// 初回時
	if nextToken == "" {
		s, _ := strconv.Atoi(size)
		n, err := gn.NotifRepo.GetByToUserOrderByLatest(
			ctx,
			gn.DB,
			userID,
			s,
			"n.id",
			"n.is_checked",
			"n.description",
			"n.created_at",
			"nt.title",
		)
		if err != nil {
			return GetNotificationsResponse{}, errors.Wrap(err, "failed to get notifications")
		}
		ns = n
	}
	if nextToken != "" {
		nt, err := strconv.Atoi(nextToken)
		if err != nil {
			return GetNotificationsResponse{}, nil
		}
		s, _ := strconv.Atoi(size)
		// お知らせ一覧を取得
		ns, err = gn.NotifRepo.GetByToUserByStartIdOrderByLatest(
			ctx,
			gn.DB,
			userID,
			model.NotificationID(nt),
			s,
			"n.id",
			"n.is_checked",
			"n.description",
			"n.created_at",
			"nt.title",
		)
		if err != nil {
			return GetNotificationsResponse{}, errors.Wrap(err, "failed to get notifications")
		}
	}

	// レスポンス作成
	res := GetNotificationsResponse{
		Notifications: []struct {
			ID          model.NotificationID
			Title       string
			Description string
			IsChecked   bool
			CreatedAt   string
		}{},
	}
	for _, n := range ns {
		res.Notifications = append(res.Notifications, struct {
			ID          model.NotificationID
			Title       string
			Description string
			IsChecked   bool
			CreatedAt   string
		}{
			ID:          model.NotificationID(n.Notification.ID),
			Title:       n.Type.Title,
			Description: n.Notification.Description,
			IsChecked:   n.Notification.IsChecked,
			CreatedAt:   model.NewTime(n.Notification.CreatedAt).Format(),
		})
	}

	// お知らせ一覧がない場合
	if len(ns) == 0 {
		res.NextToken = "0"
		return res, nil
	}

	// 次開始のお知らせIDを設定
	last := ns[len(ns)-1]
	res.NextToken = strconv.Itoa(int(last.Notification.ID) - 1)
	return res, nil
}
