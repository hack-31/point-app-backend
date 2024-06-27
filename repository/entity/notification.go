package entity

import (
	"time"

	"github.com/hack-31/point-app-backend/domain/model"
)

type NotificationID int64

type Notification struct {
	ID          NotificationID `json:"id" db:"id"`
	TypeID      int            `json:"notificationTypeId" db:"notification_type_id"`
	Title       string         `json:"title" db:"title"`
	Description string         `json:"description" db:"description"`
	ToUserID    model.UserID   `json:"toUserId" db:"to_user_id"`
	FromUserID  model.UserID   `json:"fromUserId" db:"from_user_id"`
	IsChecked   bool           `json:"isChecked" db:"is_checked"`
	CreatedAt   time.Time      `json:"createdAt" db:"created_at"`
}
type Notifications []*Notification
