package service

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/hack-31/point-app-backend/domain"
	"github.com/hack-31/point-app-backend/domain/model"
	"github.com/hack-31/point-app-backend/repository"
	"github.com/hack-31/point-app-backend/utils"
)

type DeleteUser struct {
	Connection repository.Transacter
	Cache      domain.Cache
	UserRepo   domain.UserRepo
}

func NewDeleteUser(cache domain.Cache, repo domain.UserRepo, connection repository.Transacter) *DeleteUser {
	return &DeleteUser{Cache: cache, UserRepo: repo, Connection: connection}
}

// ユーザー削除サービス
func (du *DeleteUser) DeleteUser(ctx *gin.Context, userID model.UserID) error {
	// トランザクション開始
	if err := du.Connection.Begin(ctx); err != nil {
		return err
	}
	defer func() { _ = du.Connection.Rollback() }()

	// ユーザー削除
	rows, err := du.UserRepo.DeleteUserByID(ctx, du.Connection.DB(), userID)
	if err != nil {
		return fmt.Errorf("cannot delete user: %w", err)
	}
	if rows == 0 {
		return repository.ErrNotUser
	}

	// トランザクションコミット
	if err := du.Connection.Commit(); err != nil {
		return fmt.Errorf("cannot delete user: %w", err)
	}

	// キャッシュから削除
	ownID := utils.GetUserID(ctx)
	if userID == ownID {
		if err := du.Cache.Delete(ctx, fmt.Sprint(userID)); err != nil {
			return fmt.Errorf("cannot delete in cache: %w", err)
		}
	}

	return nil
}
