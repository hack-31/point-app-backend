package service

import (
	"fmt"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hack-31/point-app-backend/constant"
	"github.com/hack-31/point-app-backend/domain"
	"github.com/hack-31/point-app-backend/domain/model"
	"github.com/hack-31/point-app-backend/domain/service"
	"github.com/hack-31/point-app-backend/myerror"
	"github.com/hack-31/point-app-backend/repository"
	utils "github.com/hack-31/point-app-backend/utils/email"
	"github.com/jmoiron/sqlx"
)

type UpdateTemporaryEmail struct {
	DB    repository.Queryer
	Cache domain.Cache
	Repo  domain.UserRepo
}

func NewUpdateTemporaryEmail(db *sqlx.DB, cache domain.Cache, rep domain.UserRepo) *UpdateTemporaryEmail {
	return &UpdateTemporaryEmail{DB: db, Cache: cache, Repo: rep}
}

// メール仮登録サービス
//
// @params
// ctx コンテキスト
// email メールアドレス
//
// @returns
// temporaryEmailId 一時保存したメールを識別するID
func (ute *UpdateTemporaryEmail) UpdateTemporaryEmail(ctx *gin.Context, email string) (string, error) {
	// ユーザードメインサービス
	userService := service.NewUserService(ute.Repo)

	// 現在利用中のメールアドレスか確認
	existMail, err := userService.ExistByEmail(ctx, &ute.DB, email)
	if err != nil {
		return "", errors.Wrap(err, "failed to check email in UpdateTemporaryEmailService.UpdateTemporaryEmail")
	}
	if existMail {
		return "", errors.Wrap(myerror.ErrAlreadyEntry, "failed to change mail address")
	}

	// キャッシュサーバーに保存するkeyの作成
	temporaryEmailID := uuid.New().String()
	confirmCode := model.NewConfirmCode().String()
	key := fmt.Sprintf("email:%s:%s", confirmCode, temporaryEmailID)
	// キャッシュサーバーへ保存
	if err = ute.Cache.Save(ctx, key, email, time.Duration(constant.ConfirmationCodeExpiration_m)); err != nil {
		return "", errors.Wrap(err, "failed to save in cache")
	}

	// メール送信
	subject := "【ポイントアプリ】本登録を完了してください"
	body := fmt.Sprintf("ポイントアプリをご利用いただきありがとうございます。\n\n確認コードは %s です。\n\nこの確認コードの有効期限は1時間です。", confirmCode)
	_, err = utils.SendMail(email, subject, body)
	if err != nil {
		return "", errors.Wrap(err, "failed to send email in UpdateTemporaryEmailService.UpdateTemporaryEmail")
	}

	return temporaryEmailID, nil
}
