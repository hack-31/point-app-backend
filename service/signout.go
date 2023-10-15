package service

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/hack-31/point-app-backend/domain"
	"github.com/hack-31/point-app-backend/utils"
)

type Signout struct {
	Cache domain.Cache
}

func NewSignout(cache domain.Cache) *Signout {
	return &Signout{Cache: cache}
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
	uid := utils.GetUserID(ctx)

	// ユーザIDをキャッシュから削除
	if err := s.Cache.Delete(ctx, fmt.Sprint(uid)); err != nil {
		return fmt.Errorf("cannot delete in cache: %w", err)
	}

	return nil
}
