package service

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/hack-31/point-app-backend/auth"
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
func (s *Signout) Signout(ctx *gin.Context) error {
	// ユーザIDの取得
	userId, _ := ctx.Get(auth.UserID)
	uid := userId.(string)

	// ユーザIDをキャッシュから削除
	if err := s.Cache.Delete(ctx, uid); err != nil {
		return fmt.Errorf("cannot delete in cache: %w", err)
	}

	return nil
}
