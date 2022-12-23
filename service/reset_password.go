package service

import (
	"context"
	"fmt"

	"github.com/hack-31/point-app-backend/domain"
	"github.com/hack-31/point-app-backend/domain/user"
	"github.com/hack-31/point-app-backend/repository"
	utils "github.com/hack-31/point-app-backend/utils/email"
)

type ResetPassword struct {
	ExecerDB  repository.Execer
	QueryerDB repository.Queryer
	Repo      domain.UserRepo
}

// パスワードリセットサービス
//
// @params
// email メールアドレス
//
// @returns
// error
func (rp *ResetPassword) ResetPassword(ctx context.Context, email string) error {
	// メール値オブジェクト作成
	mail, err := user.NewEmail(email, rp.Repo)
	if err != nil {
		return fmt.Errorf("cannot create mail object: %w", err)
	}

	// 登録可能なメールか確認
	existMail, err := mail.Exist(ctx, &rp.QueryerDB)
	if err != nil {
		return err
	}
	if !existMail {
		return fmt.Errorf("not exist email address: %w", repository.ErrNotExistEmail)
	}

	pass, err := user.NewPasswrod("")
	if err != nil {
		return fmt.Errorf("cannot create passwrod object: %w", err)
	}

	// ランダムパスワードを生成
	randomPass := pass.CreateRandomPassword()

	// パスワードハッシュ化
	hashPass, err := randomPass.CreateHash()
	if err != nil {
		return fmt.Errorf("cannot create hash passwrod: %w", err)
	}

	// パスワードを上書き
	if err := rp.Repo.UpdatePassword(ctx, rp.ExecerDB, &email, &hashPass); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// メール送信
	subject := "【ポイントアプリ】パスワード再発行完了のお知らせ"
	body := fmt.Sprintf("ポイントアプリをご利用いただきありがとうございます。\n\nポイントアプリのパスワード再設定が完了しました。\n新しいパスワードは %s です。", randomPass)
	_, err = utils.SendMail(email, subject, body)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
