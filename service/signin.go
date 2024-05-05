package service

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/hack-31/point-app-backend/domain"
	"github.com/hack-31/point-app-backend/domain/model"
	"github.com/hack-31/point-app-backend/myerror"
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
		return "", errors.Join(err, myerror.ErrNotMatchLogInfo)
	}

	// パスワードが一致するか確認
	pwd, err := model.NewPassword(password)
	if err != nil {
		return "", errors.Wrap(err, "cannot create password object")
	}
	if isMatch, _ := pwd.IsMatch(u.Password); !isMatch {
		return "", errors.Wrap(myerror.ErrNotMatchLogInfo, "no match password")
	}

	// JWTを作成
	jwt, err := s.TokenGenerator.GenerateToken(ctx, u)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate JWT")
	}

	return string(jwt), nil
}
