package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/hack-31/point-app-backend/repository/entity"
)

// お知らせ登録
// @params
// ctx context
// db dbインスタンス
// n 登録するお知らせ情報
//
// @returns
// お知らせ
func (r *Repository) CreateNotification(ctx context.Context, db Execer, n entity.Notification) (entity.Notification, error) {
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
	result, err := db.ExecContext(ctx, sql, n.ToUserID, n.FromUserID, n.TypeID, n.Description, r.Clocker.Now())
	if err != nil {
		return n, err
	}
	ID, _ := result.LastInsertId()
	n.ID = entity.NotificationID(ID)
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
func (r *Repository) GetByToUserByStartIdOrderByLatest(ctx context.Context, db Queryer, uid entity.UserID, startID entity.NotificationID, size int, columns ...string) (entity.Notifications, error) {
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
	var n entity.Notifications
	if err := db.SelectContext(ctx, &n, sql, uid, startID, size); err != nil {
		return n, err
	}
	return n, nil
}

// 指定した受信者ユーザーIDを元にお知らせ一覧を取得
func (r *Repository) GetByToUserOrderByLatest(ctx context.Context, db Queryer, uid entity.UserID, size int, columns ...string) (entity.Notifications, error) {
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
	var n entity.Notifications
	if err := db.SelectContext(ctx, &n, sql, uid, size); err != nil {
		return n, err
	}
	return n, nil
}

// お知らせ詳細取得
// @params
// ctx context
// db dbインスタンス
// uid ユーザID
// nid お知らせID
//
// @returns
// entity.Notification お知らせ
func (r *Repository) GetNotificationByID(ctx context.Context, db Queryer, uid entity.UserID, nid entity.NotificationID) (entity.Notification, error) {
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

	var n entity.Notification
	if err := db.QueryRowxContext(ctx, query, nid, uid).StructScan(&n); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return n, ErrNotFound
		}
		return n, err
	}
	return n, nil
}

// チェックしていないお知らせ総数
// @params
// ctx context
// db dbインスタンス
// ID ユーザID
//
// @returns
// お知らせ数
func (r *Repository) GetUncheckedNotificationCount(ctx context.Context, db Queryer, uid entity.UserID) (int, error) {
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
func (r *Repository) CheckNotification(ctx context.Context, db Execer, uid entity.UserID, nid entity.NotificationID) error {
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
