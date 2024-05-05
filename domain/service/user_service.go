package service

import (
	"context"
	"database/sql"

	"github.com/cockroachdb/errors"
	"github.com/hack-31/point-app-backend/domain"
	"github.com/hack-31/point-app-backend/repository"
)

type UserService struct {
	repository domain.UserRepo
}

// ユーザードメインサービス
func NewUserService(rep domain.UserRepo) *UserService {
	return &UserService{repository: rep}
}

// emailでユーザ検索
//
// @params
// db dbインスタンス
//
// @returns
// isExist true 存在, false 存在しない
func (us *UserService) ExistByEmail(ctx context.Context, db *repository.Queryer, email string) (bool, error) {
	_, err := us.repository.FindUserByEmail(ctx, *db, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
