package model

type NotificationTypeID int64

// notificationTypeID表
const (
	// ポイント送付のお知らせ
	NotificationTypeSendingPoint = 1
)

type NotificationType struct {
	ID    NotificationTypeID `json:"id" db:"id"`
	Title string             `json:"title" db:"title"`
}
