package service

import (
	"context"
	"fmt"

	"github.com/hack-31/point-app-backend/domain"
	"github.com/hack-31/point-app-backend/domain/model"
	"github.com/hack-31/point-app-backend/repository"
	"github.com/jmoiron/sqlx"
)

type Signin struct {
	DB             repository.Queryer
	Cache          domain.Cache
	Repo           domain.UserRepo
	TokenGenerator domain.TokenGenerator
}

func NewSignin(db *sqlx.DB, rep domain.UserRepo, cache domain.Cache, jwter domain.TokenGenerator) *Signin {
	return &Signin{DB: db, Cache: cache, Repo: rep, TokenGenerator: jwter}
}

// サインインサービス
//
// @params
// ctx コンテキスト
// email メール
// password パスワード
//
// @return
// JWT
func (s *Signin) Signin(ctx context.Context, email, password string) (string, error) {
	// emailよりユーザ情報を取得
	u, err := s.Repo.FindUserByEmail(ctx, s.DB, email)
	if err != nil {
		return "", fmt.Errorf("failed to find user : %w", repository.ErrNotMatchLogInfo)
	}

	// パスワードが一致するか確認
	pwd, err := model.NewPassword(password)
	if err != nil {
		return "", fmt.Errorf("cannot create password object: %w", err)
	}
	if isMatch, _ := pwd.IsMatch(u.Password); !isMatch {
		return "", fmt.Errorf("no match password:  %w", repository.ErrNotMatchLogInfo)
	}

	// JWTを作成
	jwt, err := s.TokenGenerator.GenerateToken(ctx, u)
	if err != nil {
		return "", fmt.Errorf("failed to generate JWT: %w", err)
	}

	return string(jwt), nil
}
