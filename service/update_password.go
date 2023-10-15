package service

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/hack-31/point-app-backend/domain"
	"github.com/hack-31/point-app-backend/domain/model"
	"github.com/hack-31/point-app-backend/repository"
	"github.com/hack-31/point-app-backend/utils"
	"github.com/jmoiron/sqlx"
)

type UpdatePassword struct {
	ExecerDB  repository.Execer
	QueryerDB repository.Queryer
	UserRepo  domain.UserRepo
}

func NewUpdatePassword(db *sqlx.DB, repo domain.UserRepo) *UpdatePassword {
	return &UpdatePassword{ExecerDB: db, QueryerDB: db, UserRepo: repo}
}

// パスワード更新サービス
//
// @params
// ctx コンテキスト
// oldPassword 古いパスワード
// newPassword 新しいパスワード
func (up *UpdatePassword) UpdatePassword(ctx *gin.Context, oldPassword, newPassword string) error {
	// コンテキストよりEmailを取得
	mail := utils.GetEmail(ctx)

	// Emailよりユーザ情報を取得する
	u, err := up.UserRepo.FindUserByEmail(ctx, up.QueryerDB, &mail)
	if err != nil {
		return err
	}

	// パスワードが一致するか確認
	oldPwd, err := model.NewPassword(oldPassword)
	if err != nil {
		return fmt.Errorf("cannot create password object: %w", err)
	}
	if isMatch, _ := oldPwd.IsMatch(u.Password); !isMatch {
		return fmt.Errorf("no match password:  %w", repository.ErrDifferentPassword)
	}

	// パスワード更新
	newPwd, err := model.NewPassword(newPassword)
	if err != nil {
		return fmt.Errorf("cannot create password object: %w", err)
	}
	hashNewPass, err := newPwd.CreateHash()
	if err != nil {
		return fmt.Errorf("cannot create hash password: %w", err)
	}
	if err := up.UserRepo.UpdatePassword(ctx, up.ExecerDB, &mail, &hashNewPass); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}
