package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hack-31/point-app-backend/config"
	"github.com/hack-31/point-app-backend/constant"
	"github.com/hack-31/point-app-backend/domain"
	"github.com/hack-31/point-app-backend/domain/user"
	"github.com/hack-31/point-app-backend/repository"
	utils "github.com/hack-31/point-app-backend/utils/email"
)

type RegisterTemporaryUser struct {
	DB    repository.Queryer
	Cache *repository.KVS
	Repo  domain.RegisterTemporaryUserRep
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
func (r *RegisterTemporaryUser) RegisterTemporaryUser(ctx context.Context, firstName, firstNameKana, familyName, familyNameKana, email, password string) error {
	// メール値オブジェクト作成
	mail, err := user.NewEmail(email, r.Repo)
	if err != nil {
		return fmt.Errorf("cannot create mail object: %w", err)
	}
	pass, err := user.NewPasswrod(password)
	if err != nil {
		return fmt.Errorf("cannot create passwrod object: %w", err)
	}
	// ハッシュ化
	hashPass, err := pass.CreateHash()
	if err != nil {
		return fmt.Errorf("cannot create hash passwrod: %w", err)
	}

	// 登録可能なメールか確認
	existMail, err := mail.Exist(ctx, &r.DB)
	if err != nil {
		return err
	}
	if existMail {
		return fmt.Errorf("failed to register: %w", repository.ErrAlreadyEntry)
	}

	// ユーザ情報をキャッシュに保存
	tempUserInfo := user.NewTemporaryUserString("")
	uid := uuid.New().String()
	userString := tempUserInfo.Join(firstName, firstNameKana, familyName, familyNameKana, email, hashPass)
	err = r.Cache.Save(ctx, uid, userString, time.Duration(constant.ConfirmationCodeExpiration_m))
	if err != nil {
		return fmt.Errorf("failed to save in cache: %w", err)
	}

	// メール送信
	cfg, err := config.New()
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/validation/%s", cfg.FrontEndpoint, uid)
	subject := "【ポイントアプリ】本登録を完了してください"
	body := fmt.Sprintf("%s %sさん\n\nポイントアプリをご利用いただきありがとうございます。\n\n以下のURLをクリックして、本登録を完了してください。\n%s\n\nこのURLの有効期限は1時間です。", familyName, firstName, url)
	_, err = utils.SendMail(email, subject, body)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
