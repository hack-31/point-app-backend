package service

import (
	"context"

	"github.com/hack-31/point-app-backend/domain"
	"github.com/hack-31/point-app-backend/domain/model"
	"github.com/hack-31/point-app-backend/repository"
	"github.com/jmoiron/sqlx"
)

type GetUsers struct {
	DB             repository.Queryer
	Repo           domain.UserRepo
	TokenGenerator domain.TokenGenerator
}

func NewGetUsers(db *sqlx.DB, repo domain.UserRepo, jwter domain.TokenGenerator) *GetUsers {
	return &GetUsers{DB: db, Repo: repo, TokenGenerator: jwter}
}

// ユーザ一覧取得サービス
//
// @params ctx コンテキスト
//
// @return
// ユーザ一覧
func (r *GetUsers) GetUsers(ctx context.Context) (model.Users, error) {
	// ユーザ一覧を取得する
	users, err := r.Repo.FindUsers(ctx, r.DB)
	if err != nil {
		return users, err
	}
	return users, nil
}
