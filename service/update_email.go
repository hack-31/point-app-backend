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
func (ue *UpdateEmail) UpdateEmail(ctx context.Context, temporaryEmailID, confirmCode string) {
	println("サービステスト")
	// 一時メールアドレスの復元
	key := fmt.Sprintf("email:%s:%s", confirmCode, temporaryEmailID)
	e, err := ue.Cache.Load(ctx, key)
	if err != nil {
		fmt.Println("エラーですよ！")
		return
	}
	fmt.Println(e)

}
