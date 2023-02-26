package service

import (
	"github.com/gin-gonic/gin"
	"github.com/hack-31/point-app-backend/auth"
	"github.com/hack-31/point-app-backend/domain"
	"github.com/hack-31/point-app-backend/domain/model"
	"github.com/hack-31/point-app-backend/repository"
	"github.com/jmoiron/sqlx"
)

type GetAccount struct {
	DB   repository.Queryer
	Repo domain.UserRepo
}

func NewGetAccount(db *sqlx.DB, repo domain.UserRepo) *GetAccount {
	return &GetAccount{DB: db, Repo: repo}
}

// ユーザ一覧取得サービス
//
// @params ctx コンテキスト
//
// @return
// ユーザ一覧
func (ga *GetAccount) GetAccount(ctx *gin.Context) (model.User, error) {
	// コンテキストよりEmailを取得
	email, _ := ctx.Get(auth.Email)
	stringMail := email.(string)

	// Emailよりユーザ情報を取得する
	user, err := ga.Repo.FindUserByEmail(ctx, ga.DB, &stringMail)
	if err != nil {
		return user, err
	}

	return user, nil
}
