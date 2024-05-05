package service

import (
	"fmt"

	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
	"github.com/hack-31/point-app-backend/domain"
	"github.com/hack-31/point-app-backend/domain/service"
	"github.com/hack-31/point-app-backend/myerror"
	"github.com/hack-31/point-app-backend/repository"
	"github.com/hack-31/point-app-backend/utils"
	"github.com/jmoiron/sqlx"
)

type UpdateEmail struct {
	ExecerDB  repository.Execer
	QueryerDB repository.Queryer
	Cache     domain.Cache
	UserRepo  domain.UserRepo
}

func NewUpdateEmail(db *sqlx.DB, cache domain.Cache, rep domain.UserRepo) *UpdateEmail {
	return &UpdateEmail{ExecerDB: db, QueryerDB: db, Cache: cache, UserRepo: rep}
}

// メール本変更サービス
//
// @params temporaryEmailID 一時保存メールアドレスID
// @params confirmCode 認証コード
//
// @returns
// error
func (ue *UpdateEmail) UpdateEmail(ctx *gin.Context, temporaryEmailID, confirmCode string) error {
	// ユーザードメインサービス
	userService := service.NewUserService(ue.UserRepo)
	userID := utils.GetUserID(ctx)

	// 一時メールアドレスの復元
	key := fmt.Sprintf("email:%s:%s", confirmCode, temporaryEmailID)
	temporaryEmail, err := ue.Cache.Load(ctx, key)
	if err != nil {
		return errors.Wrap(err, "failed to load email in cache")
	}

	// 復元が成功したら一時メールアドレスを削除する
	if err := ue.Cache.Delete(ctx, key); err != nil {
		return errors.Wrap(err, "failed to delete in cache")
	}

	// すでに登録済みのユーザーがいるか確認する
	existMail, err := userService.ExistByEmail(ctx, &ue.QueryerDB, temporaryEmail)
	if err != nil {
		return errors.Wrap(err, "failed to find user by mail in db")
	}
	if existMail {
		return errors.Wrap(myerror.ErrAlreadyEntry, "failed to change mail address")
	}

	// DBに保存する
	if err := ue.UserRepo.UpdateEmail(ctx, ue.ExecerDB, userID, temporaryEmail); err != nil {
		return errors.Wrap(err, "failed to update email")
	}

	// 成功時
	return nil
}
