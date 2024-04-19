package domain

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hack-31/point-app-backend/repository"
	"github.com/hack-31/point-app-backend/repository/entity"
)

// Userに対するインターフェース
type UserRepo interface {
	FindUserByEmail(ctx context.Context, db repository.Queryer, e string, columns ...string) (entity.User, error)
	GetUserByID(ctx context.Context, db repository.Queryer, ID entity.UserID) (entity.User, error)
	DeleteUserByID(ctx context.Context, db repository.Execer, ID entity.UserID) (int64, error)
	RegisterUser(ctx context.Context, db repository.Execer, u *entity.User) error
	UpdatePassword(ctx context.Context, db repository.Execer, email, pass *string) error
	UpdateEmail(ctx context.Context, db repository.Execer, userID entity.UserID, newEmail string) error
	UpdateAccount(ctx context.Context, db repository.Execer, email, familyName, familyNameKana, firstName, firstNameKana *string) error
	GetAll(ctx context.Context, db repository.Queryer, columns ...string) (entity.Users, error)
}

// 取引に対するインターフェース
type TransactionRepo interface {
	GetAquistionPoint(ctx context.Context, db repository.Queryer, userIDs []entity.UserID) (map[entity.UserID]int, error)
}

// ポイントに対するリポジトリインターフェース
type PointRepo interface {
	RegisterPointTransaction(ctx context.Context, db repository.Execer, fromUserID, toUserId entity.UserID, sendPoint int) error
	UpdateSendablePoint(ctx context.Context, db repository.Execer, fromUserID entity.UserID, sendPoint int) error
	UpdateAllSendablePoint(ctx context.Context, db repository.Execer, point int) error
}

// お知らせに対するリポジトリインターフェース
type NotificationRepo interface {
	CreateNotification(ctx context.Context, db repository.Execer, notification entity.Notification) (entity.Notification, error)
	GetByToUserByStartIdOrderByLatest(ctx context.Context, db repository.Queryer, uid entity.UserID, startID entity.NotificationID, size int, columns ...string) (entity.Notifications, error)
	GetByToUserOrderByLatest(ctx context.Context, db repository.Queryer, uid entity.UserID, size int, columns ...string) (entity.Notifications, error)
	GetNotificationByID(ctx context.Context, db repository.Queryer, uid entity.UserID, nid entity.NotificationID) (entity.Notification, error)
	GetUncheckedNotificationCount(ctx context.Context, db repository.Queryer, uid entity.UserID) (int, error)
	CheckNotification(ctx context.Context, db repository.Execer, uid entity.UserID, nid entity.NotificationID) error
}

// トークンに対するインターフェース
type TokenGenerator interface {
	GenerateToken(ctx context.Context, u entity.User) ([]byte, error)
}

// キャッシュに対するインターフェース
type Cache interface {
	Save(ctx context.Context, key, value string, minute time.Duration) error
	Load(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
	Expire(ctx context.Context, key string, minitue time.Duration) error
	Publish(ctx context.Context, channel, palyload string) error
	Subscribe(ctx *gin.Context, channel string) (<-chan string, error)
}
