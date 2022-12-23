package service

import (
	"context"
	"fmt"

	"github.com/hack-31/point-app-backend/domain"
)

type Signout struct {
	Cache domain.Cache
}

// サインアウトサービス
//
// @params
// ctx コンテキスト
// uid ユーザーID
//
// @return
// err
func (s *Signout) Signout(ctx context.Context, uid string) error {
	// ユーザーIDの存在確認
	_, err := s.Cache.Load(ctx, uid)
	if err != nil {
		return fmt.Errorf("cannot delete in cache: %w", err)
	}

	// キャッシュ削除実行
	if err := s.Cache.Delete(ctx, uid); err != nil {
		return fmt.Errorf("cannot delete in cache: %w", err)
	}

	return nil
}
