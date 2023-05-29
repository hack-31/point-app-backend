package service

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/hack-31/point-app-backend/auth"
	"github.com/hack-31/point-app-backend/domain"
	"github.com/hack-31/point-app-backend/repository"
	"github.com/jmoiron/sqlx"
)

type UpdateEmail struct {
	ExecerDB  repository.Execer
	QueryerDB repository.Queryer
	Cache     domain.Cache
	Repo      domain.UserRepo
}

func NewUpdateEmail(db *sqlx.DB, cache domain.Cache, rep domain.UserRepo) *UpdateEmail {
	return &UpdateEmail{ExecerDB: db, QueryerDB: db, Cache: cache, Repo: rep}
}

// メール本変更サービス
//
// @params temporaryEmailID 一時保存メールアドレスID
// @params confirmCode 認証コード
//
//
func (ue *UpdateEmail) UpdateEmail(ctx *gin.Context, temporaryEmailID, confirmCode string) error {
	// コンテキストよりEmailを取得する
	email, _ := ctx.Get(auth.Email)
	stringMail := email.(string)

	// 一時メールアドレスの復元
	key := fmt.Sprintf("email:%s:%s", confirmCode, temporaryEmailID)
	temporaryEmail, err := ue.Cache.Load(ctx, key)
	if err != nil {
		// TODO: エラーハンドリング
		return fmt.Errorf("cannot load user in cache: %w", err)
	}

	// 復元が成功したら一時メールアドレスを削除する
	if err := ue.Cache.Delete(ctx, key); err != nil {
		// TODO: エラーハンドリング
		return fmt.Errorf("cannot load user in cache: %w", err)
	}

	u, err := ue.Repo.FindUserByEmail(ctx, ue.QueryerDB, &stringMail)
	if err != nil {
		// TODO: エラーハンドリング
		return fmt.Errorf("not user: %w", err)
	}

	println(temporaryEmail)
	fmt.Printf("%+v", &u)

	// DBに保存する

	// TODO: 成功レスポンス
	return nil
}
