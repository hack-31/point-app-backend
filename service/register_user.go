package service

import (
	"context"
	"fmt"

	"github.com/cockroachdb/errors"
	"github.com/hack-31/point-app-backend/constant"
	"github.com/hack-31/point-app-backend/domain"
	"github.com/hack-31/point-app-backend/domain/model"
	"github.com/hack-31/point-app-backend/repository"
	"github.com/hack-31/point-app-backend/repository/entities"
	"github.com/jmoiron/sqlx"
)

type RegisterUser struct {
	DB             repository.Execer
	Cache          domain.Cache
	Repo           domain.UserRepo
	TokenGenerator domain.TokenGenerator
}

func NewRegisterUser(db *sqlx.DB, rep domain.UserRepo, cache domain.Cache, jwter domain.TokenGenerator) *RegisterUser {
	return &RegisterUser{DB: db, Cache: cache, Repo: rep, TokenGenerator: jwter}
}

// ユーザ登録サービス
//
// @params ctx コンテキスト, temporaryUserId 一時保存ユーザid
//
// @return ユーザ情報
func (r *RegisterUser) RegisterUser(ctx context.Context, temporaryUserId, confirmCode string) (*entities.User, string, error) {
	// 一時ユーザ情報を復元
	key := fmt.Sprintf("user:%s:%s", confirmCode, temporaryUserId)
	u, err := r.Cache.Load(ctx, key)
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to load in cache")
	}

	// 復元が成功したら一時ユーザ情報除削
	if err := r.Cache.Delete(ctx, key); err != nil {
		return nil, "", errors.Wrap(err, "failed to delete in cache")
	}

	// 復元したユーザ情報を解析
	temporyUser := model.NewTemporaryUserString(u)
	firstName, firstNameKana, familyName, familyNameKana, email, hashPass := temporyUser.Split()

	// DBに保存
	user := &entities.User{
		FirstName:      firstName,
		FirstNameKana:  firstNameKana,
		FamilyName:     familyName,
		FamilyNameKana: familyNameKana,
		Email:          email,
		Password:       hashPass,
		SendingPoint:   constant.DefaultSendingPoint,
	}
	if err := r.Repo.RegisterUser(ctx, r.DB, user); err != nil {
		return nil, "", errors.Wrap(err, "failed to register")
	}

	// JWTを作成
	jwt, err := r.TokenGenerator.GenerateToken(ctx, *user)
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to generate JWT")
	}

	return user, string(jwt), nil
}
