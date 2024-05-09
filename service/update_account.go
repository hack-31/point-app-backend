package service

import (
	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
	"github.com/hack-31/point-app-backend/domain"
	"github.com/hack-31/point-app-backend/repository"
	"github.com/hack-31/point-app-backend/utils"
	"github.com/jmoiron/sqlx"
)

type UpdateAccount struct {
	ExecerDB repository.Execer
	UserRepo domain.UserRepo
}

func NewUpdateAccount(db *sqlx.DB, repo domain.UserRepo) *UpdateAccount {
	return &UpdateAccount{ExecerDB: db, UserRepo: repo}
}

// アカウント情報更新サービス
//
// @params
// ctx コンテキスト
// familyName 苗字
// familyNameKana 苗字カナ
// firstName 名前
// firstNameKana 名前カナ
func (ua *UpdateAccount) UpdateAccount(ctx *gin.Context, familyName, familyNameKana, firstName, firstNameKana string) error {
	// コンテキストよりEmailを取得

	mail := utils.GetEmail(ctx)
	// アカウント情報更新
	if err := ua.UserRepo.UpdateAccount(ctx, ua.ExecerDB, &mail, &familyName, &familyNameKana, &firstName, &firstNameKana); err != nil {
		return errors.Wrap(err, "failed to update account")
	}

	return nil
}
