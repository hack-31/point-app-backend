package service

import (
	"fmt"

	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
	"github.com/hack-31/point-app-backend/domain"
	"github.com/hack-31/point-app-backend/myerror"
	"github.com/hack-31/point-app-backend/repository"
	"github.com/hack-31/point-app-backend/repository/entity"
	"github.com/hack-31/point-app-backend/utils"
	"github.com/jmoiron/sqlx"
)

type DeleteUser struct {
	Tx       repository.Beginner
	Cache    domain.Cache
	UserRepo domain.UserRepo
}

func NewDeleteUser(cache domain.Cache, repo domain.UserRepo, db *sqlx.DB) *DeleteUser {
	return &DeleteUser{Cache: cache, UserRepo: repo, Tx: db}
}

// ユーザー削除サービス
func (du *DeleteUser) DeleteUser(ctx *gin.Context, userID entity.UserID) error {
	// トランザクション開始
	tx, err := du.Tx.BeginTxx(ctx, nil)
	defer func() { _ = tx.Rollback() }()
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}

	// ユーザー削除
	rows, err := du.UserRepo.DeleteUserByID(ctx, tx, userID)
	if err != nil {
		return errors.Wrap(err, "failed to delete user")
	}
	if rows == 0 {
		return errors.Wrap(myerror.ErrNotUser, "failed to delete user")
	}

	// トランザクションコミット
	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	// キャッシュから削除
	ownID := utils.GetUserID(ctx)
	if userID == ownID {
		if err := du.Cache.Delete(ctx, fmt.Sprint(userID)); err != nil {
			return errors.Wrap(err, "failed to delete in cache")
		}
	}

	return nil
}
