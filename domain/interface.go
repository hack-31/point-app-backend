package domain

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hack-31/point-app-backend/domain/model"
	"github.com/hack-31/point-app-backend/repository"
	customentities "github.com/hack-31/point-app-backend/repository/custom_entities"
	"github.com/hack-31/point-app-backend/repository/entities"
)

// Userに対するインターフェース
type UserRepo interface {
	FindUserByEmail(ctx context.Context, db repository.Queryer, e string, columns ...string) (entities.User, error)
	GetUserByID(ctx context.Context, db repository.Queryer, ID model.UserID) (entities.User, error)
	DeleteUserByID(ctx context.Context, db repository.Execer, ID model.UserID) (int64, error)
	RegisterUser(ctx context.Context, db repository.Execer, u *entities.User) error
	UpdatePassword(ctx context.Context, db repository.Execer, email, pass *string) error
	UpdateEmail(ctx context.Context, db repository.Execer, userID model.UserID, newEmail string) error
	UpdateAccount(ctx context.Context, db repository.Execer, email, familyName, familyNameKana, firstName, firstNameKana *string) error
	GetAll(ctx context.Context, db repository.Queryer, columns ...string) ([]*entities.User, error)
	GetAllWithCursor(ctx context.Context, db repository.Queryer, param repository.GetAllWithCursorParam) ([]*entities.User, error)
	GetUsers(ctx context.Context, db repository.Queryer, param repository.GetUsersParam) ([]*entities.User, error)
}

// 取引に対するインターフェース
type TransactionRepo interface {
	GetAquistionPoint(ctx context.Context, db repository.Queryer, userIDs []model.UserID) (map[model.UserID]int, error)
}

// ポイントに対するリポジトリインターフェース
type PointRepo interface {
	RegisterPointTransaction(ctx context.Context, db repository.Execer, fromUserID, toUserId model.UserID, sendPoint int) error
	UpdateSendablePoint(ctx context.Context, db repository.Execer, fromUserID model.UserID, sendPoint int) error
	UpdateAllSendablePoint(ctx context.Context, db repository.Execer, point int) error
}

// お知らせに対するリポジトリインターフェース
type NotificationRepo interface {
	CreateNotification(ctx context.Context, db repository.Execer, notification customentities.Notification) (customentities.Notification, error)
	GetByToUserByStartIdOrderByLatest(ctx context.Context, db repository.Queryer, uid model.UserID, startID model.NotificationID, size int, columns ...string) ([]*customentities.Notification, error)
	GetByToUserOrderByLatest(ctx context.Context, db repository.Queryer, uid model.UserID, size int, columns ...string) ([]*customentities.Notification, error)
	GetNotificationByID(ctx context.Context, db repository.Queryer, uid model.UserID, nid model.NotificationID) (customentities.Notification, error)
	GetUncheckedNotificationCount(ctx context.Context, db repository.Queryer, uid model.UserID) (int, error)
	CheckNotification(ctx context.Context, db repository.Execer, uid model.UserID, nid model.NotificationID) error
}

// トークンに対するインターフェース
type TokenGenerator interface {
	GenerateToken(ctx context.Context, u entities.User) ([]byte, error)
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
