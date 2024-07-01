package customentities

import "github.com/hack-31/point-app-backend/repository/entities"

// NotificationTypeとNotificationを組み合わせたカスタムエンティティ
type Notification struct {
	Notification entities.Notification
	Type         entities.NotificationType
}
