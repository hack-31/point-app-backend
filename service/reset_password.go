package service

import (
	"context"
	"fmt"

	"github.com/cockroachdb/errors"
	"github.com/hack-31/point-app-backend/domain"
	"github.com/hack-31/point-app-backend/domain/model"
	"github.com/hack-31/point-app-backend/domain/service"
	"github.com/hack-31/point-app-backend/myerror"
	"github.com/hack-31/point-app-backend/repository"
	utils "github.com/hack-31/point-app-backend/utils/email"
	"github.com/jmoiron/sqlx"
)

type ResetPassword struct {
	ExecerDB  repository.Execer
	QueryerDB repository.Queryer
	Repo      domain.UserRepo
}

func NewResetPassword(db *sqlx.DB, rep domain.UserRepo) *ResetPassword {
	return &ResetPassword{ExecerDB: db, QueryerDB: db, Repo: rep}
}

// パスワードリセットサービス
//
// @params
// email メールアドレス
//
// @returns
// error
func (rp *ResetPassword) ResetPassword(ctx context.Context, email string) error {
	// ユーザドメインサービス
	userService := service.NewUserService(rp.Repo)

	// 登録可能なメールか確認
	existMail, err := userService.ExistByEmail(ctx, &rp.QueryerDB, email)
	if err != nil {
		return errors.Wrap(err, "failed to check email in ResetPasswordService.ResetPassword")
	}
	if !existMail {
		return errors.Wrap(myerror.ErrNotExistEmail, "failed to reset password in ResetPasswordService.ResetPassword")
	}

	pass, err := model.NewPassword("")
	if err != nil {
		return errors.Wrap(err, "cannot create password object")
	}

	// ランダムパスワードを生成
	randomPass := pass.CreateRandomPassword()

	// パスワードハッシュ化
	hashPass, err := randomPass.CreateHash()
	if err != nil {
		return errors.Wrap(err, "failed to create hash password in ResetPasswordService.ResetPassword")
	}

	// パスワードを上書き
	if err := rp.Repo.UpdatePassword(ctx, rp.ExecerDB, &email, &hashPass); err != nil {
		return errors.Wrap(err, "failed to update password in ResetPasswordService.ResetPassword")
	}

	// メール送信
	subject := "【ポイントアプリ】パスワード再発行完了のお知らせ"
	body := fmt.Sprintf("ポイントアプリをご利用いただきありがとうございます。\n\nポイントアプリのパスワード再設定が完了しました。\n新しいパスワードは %s です。", randomPass.String())
	_, err = utils.SendMail(email, subject, body)
	if err != nil {
		return errors.Wrap(err, "failed to send email in ResetPasswordService.ResetPassword")
	}

	return nil
}
