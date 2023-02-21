package service

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/hack-31/point-app-backend/auth"
	"github.com/hack-31/point-app-backend/domain"
	"github.com/hack-31/point-app-backend/repository"
)

type UpdateAccount struct {
	ExecerDB  repository.Execer
	QueryerDB repository.Queryer
	UserRepo  domain.UserRepo
}

// アカウント情報更新サービス
//
// @params
// ctx コンテキスト
// familyName 苗字
// familNameKana 苗字カナ
// firstName 名前
// firstNameKana 名前カナ
func (ua *UpdateAccount) UpdateAccount(ctx *gin.Context, familyName, familyNameKana, firstName, firstNameKana string) error {
	// コンテキストよりEmailを取得
	email, _ := ctx.Get(auth.Email)
	stringMail := email.(string)

	// Emailよりユーザ情報を取得する
	_, err := ua.UserRepo.FindUserByEmail(ctx, ua.QueryerDB, &stringMail)
	if err != nil {
		return fmt.Errorf("no match user: %w", repository.ErrNotUser)
	}

	// アカウント情報更新
	if err := ua.UserRepo.UpdateAccount(ctx, ua.ExecerDB, &stringMail, &familyName, &familyNameKana, &firstName, &firstNameKana); err != nil {
		return fmt.Errorf("failed to update account: %w", err)
	}

	return nil
}
