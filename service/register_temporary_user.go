package service

import (
	"context"
	"fmt"
	"time"

	"github.com/cockroachdb/errors"
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

type RegisterTemporaryUser struct {
	DB    repository.Queryer
	Cache domain.Cache
	Repo  domain.UserRepo
}

func NewRegisterTemporaryUser(db *sqlx.DB, rep domain.UserRepo, cache domain.Cache) *RegisterTemporaryUser {
	return &RegisterTemporaryUser{DB: db, Cache: cache, Repo: rep}
}

// ユーザ仮登録サービス
//
// @params
// ctx コンテキスト
// firstName 名前
// firstNameKana 名前カナ
// familyName 名字
// familyNameKana 名字カナ
// email メールアドレス
// password パスワード
//
// @returns
// temporaryUserId 一時保存したユーザを識別するID
func (r *RegisterTemporaryUser) RegisterTemporaryUser(ctx context.Context, firstName, firstNameKana, familyName, familyNameKana, email, password string) (string, error) {
	// ユーザドメインサービス
	userService := service.NewUserService(r.Repo)

	// 登録可能なメールか確認
	existMail, err := userService.ExistByEmail(ctx, &r.DB, email)
	if err != nil {
		return "", errors.Wrap(err, "failed to check email in RegisterTemporaryUserService.RegisterTemporaryUser")
	}
	if existMail {
		return "", errors.Wrap(myerror.ErrAlreadyEntry, "failed to register in RegisterTemporaryUserService.RegisterTemporaryUser")
	}

	// パスワードハッシュ化
	pass, err := model.NewPassword(password)
	if err != nil {
		return "", errors.Wrap(err, "failed to create password object in RegisterTemporaryUserService.RegisterTemporaryUser")
	}
	hashPass, err := pass.CreateHash()
	if err != nil {
		return "", errors.Wrap(err, "failed to create hash password in RegisterTemporaryUserService.RegisterTemporaryUser")
	}

	// ユーザ情報をキャッシュに保存
	tempUserInfo := model.NewTemporaryUserString("")
	// キャッシュサーバーに保存するkeyの作成
	uid := uuid.New().String()
	confirmCode := model.NewConfirmCode().String()
	key := fmt.Sprintf("user:%s:%s", confirmCode, uid)
	// キャッシュのサーバーに保存するvalueを作成
	userString := tempUserInfo.Join(firstName, firstNameKana, familyName, familyNameKana, email, hashPass)
	// 保存
	err = r.Cache.Save(ctx, key, userString, time.Duration(constant.ConfirmationCodeExpiration_m))
	if err != nil {
		return "", errors.Wrap(err, "failed to save in cache in RegisterTemporaryUserService.RegisterTemporaryUser")
	}

	// メール送信
	subject := "【ポイントアプリ】本登録を完了してください"
	body := fmt.Sprintf("%s %sさん\n\nポイントアプリをご利用いただきありがとうございます。\n\n確認コードは %s です。\n\nこの確認コードの有効期限は1時間です。", familyName, firstName, confirmCode)
	_, err = utils.SendMail(email, subject, body)
	if err != nil {
		return "", errors.Wrap(err, "failed to send email in RegisterTemporaryUserService.RegisterTemporaryUser")
	}

	return uid, nil
}
