package service

import (
	"context"
	"fmt"

	"github.com/hack-31/point-app-backend/domain"
)

type UpdateEmail struct {
	Cache domain.Cache
}

func NewUpdateEmail(cache domain.Cache) *UpdateEmail {
	return &UpdateEmail{Cache: cache}
}

// メール本変更サービス
//
// @params temporaryEmailID 一時保存メールアドレスID
// @params confirmCode 認証コード
//
//
func (ue *UpdateEmail) UpdateEmail(ctx context.Context, temporaryEmailID, confirmCode string) error {
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

	println(temporaryEmail)

	// DBに保存する

	// TODO: 成功レスポンス
	return nil
}
