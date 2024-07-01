package repository

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/hack-31/point-app-backend/domain/model"
	"github.com/hack-31/point-app-backend/myerror"
	customentities "github.com/hack-31/point-app-backend/repository/custom_entities"
	"github.com/hack-31/point-app-backend/repository/entities"
)

// お知らせ登録
// @params
// ctx context
// db dbインスタンス
// n 登録するお知らせ情報
//
// @returns
// お知らせ
func (r *Repository) CreateNotification(ctx context.Context, db Execer, n customentities.Notification) (customentities.Notification, error) {
	sql := `
		INSERT INTO notifications (
			to_user_id,
			from_user_id,
			notification_type_id,
			description,
			created_at		
		)
		VALUES (?,?,?,?,?)
	`
	result, err := db.ExecContext(ctx, sql, n.Notification.ToUserID, n.Notification.FromUserID, n.Type.ID, n.Notification.Description, r.Clocker.Now())
	if err != nil {
		return n, err
	}
	ID, _ := result.LastInsertId()
	n.Notification.ID = uint64(ID)
	return n, nil
}

// お知らせ一覧取得
// @params
// ctx context
// db dbインスタンス
// uid 受信者ユーザID
// startID 開始するお知らせID
// size 取得件数
//
// @returns
// お知らせ一覧
func (r *Repository) GetByToUserByStartIdOrderByLatest(ctx context.Context, db Queryer, uid model.UserID, startID model.NotificationID, size int, columns ...string) ([]*customentities.Notification, error) {
	column := "n.id, n.is_checked, n.to_user_id, n.from_user_id, n.description, n.created_at, nt.title"
	if len(columns) > 0 {
		column = strings.Join(columns, ", ")
	}

	sql := "SELECT " + column + " " +
		"FROM notifications AS n " +
		"LEFT JOIN notification_types AS nt " +
		"ON nt.id = n.notification_type_id " +
		"WHERE n.to_user_id = ? AND n.id <= ? " +
		"ORDER BY n.id DESC " +
		"LIMIT ?;"

	var n []*struct {
		ID          uint64       `json:"id" db:"id"`
		TypeID      int          `json:"notificationTypeId" db:"notification_type_id"`
		Title       string       `json:"title" db:"title"`
		Description string       `json:"description" db:"description"`
		ToUserID    model.UserID `json:"toUserId" db:"to_user_id"`
		FromUserID  model.UserID `json:"fromUserId" db:"from_user_id"`
		IsChecked   bool         `json:"isChecked" db:"is_checked"`
		CreatedAt   time.Time    `json:"createdAt" db:"created_at"`
	}
	if err := db.SelectContext(ctx, &n, sql, uid, startID, size); err != nil {
		return []*customentities.Notification{}, err
	}
	var res []*customentities.Notification
	for _, v := range n {
		res = append(res, &customentities.Notification{
			Notification: entities.Notification{
				ID:                 v.ID,
				Description:        v.Description,
				NotificationTypeID: uint64(v.TypeID),
				ToUserID:           uint64(v.ToUserID),
				FromUserID:         uint64(v.FromUserID),
				IsChecked:          v.IsChecked,
				CreatedAt:          v.CreatedAt,
			},
			Type: entities.NotificationType{
				Title: v.Title,
			},
		})
	}
	return res, nil
}

// 指定した受信者ユーザーIDを元にお知らせ一覧を取得
func (r *Repository) GetByToUserOrderByLatest(ctx context.Context, db Queryer, uid model.UserID, size int, columns ...string) ([]*customentities.Notification, error) {
	column := "n.id, n.is_checked, n.to_user_id, n.from_user_id, n.description, n.created_at, nt.title"
	if len(columns) > 0 {
		column = strings.Join(columns, ", ")
	}

	sql := "SELECT " + column + " " +
		"FROM notifications AS n " +
		"LEFT JOIN notification_types AS nt " +
		"ON nt.id = n.notification_type_id " +
		"WHERE n.to_user_id = ? " +
		"ORDER BY n.id DESC " +
		"LIMIT ?;"

	var n []*struct {
		ID          uint64       `json:"id" db:"id"`
		TypeID      int          `json:"notificationTypeId" db:"notification_type_id"`
		Title       string       `json:"title" db:"title"`
		Description string       `json:"description" db:"description"`
		ToUserID    model.UserID `json:"toUserId" db:"to_user_id"`
		FromUserID  model.UserID `json:"fromUserId" db:"from_user_id"`
		IsChecked   bool         `json:"isChecked" db:"is_checked"`
		CreatedAt   time.Time    `json:"createdAt" db:"created_at"`
	}
	if err := db.SelectContext(ctx, &n, sql, uid, size); err != nil {
		return []*customentities.Notification{}, err
	}
	var res []*customentities.Notification
	for _, v := range n {
		res = append(res, &customentities.Notification{
			Notification: entities.Notification{
				ID:          v.ID,
				Description: v.Description,
				ToUserID:    uint64(v.ToUserID),
				FromUserID:  uint64(v.FromUserID),
				IsChecked:   v.IsChecked,
				CreatedAt:   v.CreatedAt,
			},
			Type: entities.NotificationType{
				Title: v.Title,
			},
		})
	}
	return res, nil
}

// お知らせ詳細取得
// @params
// ctx context
// db dbインスタンス
// uid ユーザID
// nid お知らせID
//
// @returns
// entities.Notification お知らせ
func (r *Repository) GetNotificationByID(ctx context.Context, db Queryer, uid model.UserID, nid model.NotificationID) (customentities.Notification, error) {
	query := `
		SELECT 
			n.id,
			n.from_user_id,
			n.from_user_id,
			n.description,
			n.created_at,
			n.is_checked,
			nt.title
		FROM notifications AS n
		LEFT JOIN notification_types AS nt
		ON nt.id = n.notification_type_id
		WHERE n.id = ? AND n.to_user_id = ?
		LIMIT 1
	`

	var n struct {
		ID          uint64       `json:"id" db:"id"`
		TypeID      int          `json:"notificationTypeId" db:"notification_type_id"`
		Title       string       `json:"title" db:"title"`
		Description string       `json:"description" db:"description"`
		ToUserID    model.UserID `json:"toUserId" db:"to_user_id"`
		FromUserID  model.UserID `json:"fromUserId" db:"from_user_id"`
		IsChecked   bool         `json:"isChecked" db:"is_checked"`
		CreatedAt   time.Time    `json:"createdAt" db:"created_at"`
	}
	var res customentities.Notification
	if err := db.QueryRowxContext(ctx, query, nid, uid).StructScan(&n); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return res, myerror.ErrNotFound
		}
		return res, err
	}

	return customentities.Notification{
		Notification: entities.Notification{
			ID:                 n.ID,
			Description:        n.Description,
			NotificationTypeID: uint64(n.TypeID),
			ToUserID:           uint64(n.ToUserID),
			FromUserID:         uint64(n.FromUserID),
			IsChecked:          n.IsChecked,
			CreatedAt:          n.CreatedAt,
		},
		Type: entities.NotificationType{
			Title: n.Title,
		},
	}, nil
}

// チェックしていないお知らせ総数
// @params
// ctx context
// db dbインスタンス
// ID ユーザID
//
// @returns
// お知らせ数
func (r *Repository) GetUncheckedNotificationCount(ctx context.Context, db Queryer, uid model.UserID) (int, error) {
	query := `
		SELECT COUNT(1) AS count
		FROM notifications
		WHERE to_user_id = ? AND is_checked = false
	`

	var cnt int
	if err := db.GetContext(ctx, &cnt, query, uid); err != nil {
		// お知らせが無いからと言ってエラーでは無い
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}
	return cnt, nil
}

// お知らせチェックフラグをONにする
// @params
// ctx context
// db dbインスタンス
// uid ユーザID
// nie お知らせID
func (r *Repository) CheckNotification(ctx context.Context, db Execer, uid model.UserID, nid model.NotificationID) error {
	query := `
		UPDATE notifications
		SET is_checked = true
		WHERE id = ? AND to_user_id = ?
	`
	_, err := db.ExecContext(ctx, query, nid, uid)
	if err != nil {
		return err
	}
	return nil
}
