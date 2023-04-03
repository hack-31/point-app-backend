package service

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/hack-31/point-app-backend/domain"
	"github.com/hack-31/point-app-backend/domain/service"
	"github.com/hack-31/point-app-backend/repository"
	"github.com/jmoiron/sqlx"
)

type RegisterTemporaryEmail struct {
	DB   repository.Queryer
	Repo domain.UserRepo
}

func NewRegisterTemporaryEmail(db *sqlx.DB, rep domain.UserRepo) *RegisterTemporaryEmail {
	return &RegisterTemporaryEmail{DB: db, Repo: rep}
}

// メール仮登録サービス
//
// @params
// ctx コンテキスト
// email メールアドレス
//
// @returns
// temporaryEmailId 一時保存したメールを識別するID
func (r *RegisterTemporaryEmail) RegisterTemporaryEmail(ctx *gin.Context, email string) (string, error) {
	// ユーザードメインサービス
	userService := service.NewUserService(r.Repo)

	// 現在利用中のメールアドレスか確認
	existMail, err := userService.ExistByEmail(ctx, &r.DB, email)
	if err != nil {
		return "", err
	}
	if existMail {
		return "", fmt.Errorf("failed to register: %w", repository.ErrAlreadyEntry)
	}
	return "成功です。", nil
}
