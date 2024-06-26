package entities

// Code generated by xo. DO NOT EDIT.

import (
	"context"
	"time"
)

// Notification represents a row from 'point_app.notifications'.
type Notification struct {
	ID                 uint64    `json:"id" db:"id"`                                     // お知らせID
	ToUserID           uint64    `json:"to_user_id" db:"to_user_id"`                     // お知らせ先ユーザID
	FromUserID         uint64    `json:"from_user_id" db:"from_user_id"`                 // お知らせ元ユーザID
	IsChecked          bool      `json:"is_checked" db:"is_checked"`                     // チェックフラグ
	NotificationTypeID uint64    `json:"notification_type_id" db:"notification_type_id"` // お知らせ種別ID
	Description        string    `json:"description" db:"description"`                   // 説明
	CreatedAt          time.Time `json:"created_at" db:"created_at"`                     // 作成日時
	// xo fields
	_exists, _deleted bool
}

// Exists returns true when the [Notification] exists in the database.
func (n *Notification) Exists() bool {
	return n._exists
}

// Deleted returns true when the [Notification] has been marked for deletion
// from the database.
func (n *Notification) Deleted() bool {
	return n._deleted
}

// Insert inserts the [Notification] to the database.
func (n *Notification) Insert(ctx context.Context, db DB) error {
	switch {
	case n._exists: // already exists
		return logerror(&ErrInsertFailed{ErrAlreadyExists})
	case n._deleted: // deleted
		return logerror(&ErrInsertFailed{ErrMarkedForDeletion})
	}
	// insert (primary key generated and returned by database)
	const sqlstr = `INSERT INTO point_app.notifications (` +
		`to_user_id, from_user_id, is_checked, notification_type_id, description, created_at` +
		`) VALUES (` +
		`?, ?, ?, ?, ?, ?` +
		`)`
	// run
	logf(sqlstr, n.ToUserID, n.FromUserID, n.IsChecked, n.NotificationTypeID, n.Description, n.CreatedAt)
	res, err := db.ExecContext(ctx, sqlstr, n.ToUserID, n.FromUserID, n.IsChecked, n.NotificationTypeID, n.Description, n.CreatedAt)
	if err != nil {
		return logerror(err)
	}
	// retrieve id
	id, err := res.LastInsertId()
	if err != nil {
		return logerror(err)
	} // set primary key
	n.ID = uint64(id)
	// set exists
	n._exists = true
	return nil
}

// Update updates a [Notification] in the database.
func (n *Notification) Update(ctx context.Context, db DB) error {
	switch {
	case !n._exists: // doesn't exist
		return logerror(&ErrUpdateFailed{ErrDoesNotExist})
	case n._deleted: // deleted
		return logerror(&ErrUpdateFailed{ErrMarkedForDeletion})
	}
	// update with primary key
	const sqlstr = `UPDATE point_app.notifications SET ` +
		`to_user_id = ?, from_user_id = ?, is_checked = ?, notification_type_id = ?, description = ?, created_at = ? ` +
		`WHERE id = ?`
	// run
	logf(sqlstr, n.ToUserID, n.FromUserID, n.IsChecked, n.NotificationTypeID, n.Description, n.CreatedAt, n.ID)
	if _, err := db.ExecContext(ctx, sqlstr, n.ToUserID, n.FromUserID, n.IsChecked, n.NotificationTypeID, n.Description, n.CreatedAt, n.ID); err != nil {
		return logerror(err)
	}
	return nil
}

// Save saves the [Notification] to the database.
func (n *Notification) Save(ctx context.Context, db DB) error {
	if n.Exists() {
		return n.Update(ctx, db)
	}
	return n.Insert(ctx, db)
}

// Upsert performs an upsert for [Notification].
func (n *Notification) Upsert(ctx context.Context, db DB) error {
	switch {
	case n._deleted: // deleted
		return logerror(&ErrUpsertFailed{ErrMarkedForDeletion})
	}
	// upsert
	const sqlstr = `INSERT INTO point_app.notifications (` +
		`id, to_user_id, from_user_id, is_checked, notification_type_id, description, created_at` +
		`) VALUES (` +
		`?, ?, ?, ?, ?, ?, ?` +
		`)` +
		` ON DUPLICATE KEY UPDATE ` +
		`to_user_id = VALUES(to_user_id), from_user_id = VALUES(from_user_id), is_checked = VALUES(is_checked), notification_type_id = VALUES(notification_type_id), description = VALUES(description), created_at = VALUES(created_at)`
	// run
	logf(sqlstr, n.ID, n.ToUserID, n.FromUserID, n.IsChecked, n.NotificationTypeID, n.Description, n.CreatedAt)
	if _, err := db.ExecContext(ctx, sqlstr, n.ID, n.ToUserID, n.FromUserID, n.IsChecked, n.NotificationTypeID, n.Description, n.CreatedAt); err != nil {
		return logerror(err)
	}
	// set exists
	n._exists = true
	return nil
}

// Delete deletes the [Notification] from the database.
func (n *Notification) Delete(ctx context.Context, db DB) error {
	switch {
	case !n._exists: // doesn't exist
		return nil
	case n._deleted: // deleted
		return nil
	}
	// delete with single primary key
	const sqlstr = `DELETE FROM point_app.notifications ` +
		`WHERE id = ?`
	// run
	logf(sqlstr, n.ID)
	if _, err := db.ExecContext(ctx, sqlstr, n.ID); err != nil {
		return logerror(err)
	}
	// set deleted
	n._deleted = true
	return nil
}

// NotificationsByNotificationTypeID retrieves a row from 'point_app.notifications' as a [Notification].
//
// Generated from index 'fk_notification_type_id'.
func NotificationsByNotificationTypeID(ctx context.Context, db DB, notificationTypeID uint64) ([]*Notification, error) {
	// query
	const sqlstr = `SELECT ` +
		`id, to_user_id, from_user_id, is_checked, notification_type_id, description, created_at ` +
		`FROM point_app.notifications ` +
		`WHERE notification_type_id = ?`
	// run
	logf(sqlstr, notificationTypeID)
	rows, err := db.QueryContext(ctx, sqlstr, notificationTypeID)
	if err != nil {
		return nil, logerror(err)
	}
	defer rows.Close()
	// process
	var res []*Notification
	for rows.Next() {
		n := Notification{
			_exists: true,
		}
		// scan
		if err := rows.Scan(&n.ID, &n.ToUserID, &n.FromUserID, &n.IsChecked, &n.NotificationTypeID, &n.Description, &n.CreatedAt); err != nil {
			return nil, logerror(err)
		}
		res = append(res, &n)
	}
	if err := rows.Err(); err != nil {
		return nil, logerror(err)
	}
	return res, nil
}

// NotificationsByToUserIDID retrieves a row from 'point_app.notifications' as a [Notification].
//
// Generated from index 'idx_to_user_id'.
func NotificationsByToUserIDID(ctx context.Context, db DB, toUserID, id uint64) ([]*Notification, error) {
	// query
	const sqlstr = `SELECT ` +
		`id, to_user_id, from_user_id, is_checked, notification_type_id, description, created_at ` +
		`FROM point_app.notifications ` +
		`WHERE to_user_id = ? AND id = ?`
	// run
	logf(sqlstr, toUserID, id)
	rows, err := db.QueryContext(ctx, sqlstr, toUserID, id)
	if err != nil {
		return nil, logerror(err)
	}
	defer rows.Close()
	// process
	var res []*Notification
	for rows.Next() {
		n := Notification{
			_exists: true,
		}
		// scan
		if err := rows.Scan(&n.ID, &n.ToUserID, &n.FromUserID, &n.IsChecked, &n.NotificationTypeID, &n.Description, &n.CreatedAt); err != nil {
			return nil, logerror(err)
		}
		res = append(res, &n)
	}
	if err := rows.Err(); err != nil {
		return nil, logerror(err)
	}
	return res, nil
}

// NotificationByID retrieves a row from 'point_app.notifications' as a [Notification].
//
// Generated from index 'notifications_id_pkey'.
func NotificationByID(ctx context.Context, db DB, id uint64) (*Notification, error) {
	// query
	const sqlstr = `SELECT ` +
		`id, to_user_id, from_user_id, is_checked, notification_type_id, description, created_at ` +
		`FROM point_app.notifications ` +
		`WHERE id = ?`
	// run
	logf(sqlstr, id)
	n := Notification{
		_exists: true,
	}
	if err := db.QueryRowContext(ctx, sqlstr, id).Scan(&n.ID, &n.ToUserID, &n.FromUserID, &n.IsChecked, &n.NotificationTypeID, &n.Description, &n.CreatedAt); err != nil {
		return nil, logerror(err)
	}
	return &n, nil
}

// NotificationType returns the NotificationType associated with the [Notification]'s (NotificationTypeID).
//
// Generated from foreign key 'fk_notification_type_id'.
func (n *Notification) NotificationType(ctx context.Context, db DB) (*NotificationType, error) {
	return NotificationTypeByID(ctx, db, n.NotificationTypeID)
}
