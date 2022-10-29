package service

import (
	"context"
	"fmt"

	"github.com/hack-31/point-app-backend/domain/user"
	"github.com/hack-31/point-app-backend/entity"
	"github.com/hack-31/point-app-backend/repository"
)

type RegisterUser struct {
	DB   repository.Execer
	Repo UserRegister
}

// ユーザ登録サービス
//
// @params ctx コンテキスト
//
// @params namae ユーザ名
//
// @param password パスワード
//
// @params email メールアドレス
//
// @params role ロール
//
// @return ユーザ情報
func (r *RegisterUser) RegisterUser(ctx context.Context, name, password, email string, role int) (*entity.User, error) {

	pwd, err := user.NewPasswrod(&password)
	if err != nil {
		return nil, fmt.Errorf("cannot password: %w", err)
	}

	hashPwd, err := pwd.CreateHash()
	if err != nil {
		return nil, fmt.Errorf("cannot hash password: %w", err)
	}

	user := &entity.User{
		Name:     name,
		Password: string(hashPwd),
		Role:     role,
		Email:    email,
	}

	if err := r.Repo.RegisterUser(ctx, r.DB, user); err != nil {
		return nil, fmt.Errorf("failed to register: %w", err)
	}
	return user, nil
}
