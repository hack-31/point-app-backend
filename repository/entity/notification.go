package entity

import (
	"time"
)

type NotificationID int64

type Notification struct {
	ID          NotificationID `json:"id" db:"id"`
	TypeID      int            `json:"notificationTypeId" db:"notification_type_id"`
	Title       string         `json:"title" db:"title"`
	Description string         `json:"description" db:"description"`
	ToUserID    UserID         `json:"toUserId" db:"to_user_id"`
	FromUserID  UserID         `json:"fromUserId" db:"from_user_id"`
	IsChecked   bool           `json:"isChecked" db:"is_checked"`
	CreatedAt   time.Time      `json:"createdAt" db:"created_at"`
}
type Notifications []*Notification
