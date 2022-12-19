package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"unicode/utf8"

	"github.com/hack-31/point-app-backend/domain"
	"github.com/hack-31/point-app-backend/repository"
)

type Email struct {
	value      string
	repository domain.UserRepo
}

func NewEmail(mail string, rep domain.UserRepo) (*Email, error) {
	if 256 < utf8.RuneCountInString(mail) {
		return nil, fmt.Errorf("cannot use email over 257 char")
	}
	return &Email{value: mail, repository: rep}, nil
}

// emailでユーザ検索
//
// @params
// db dbインスタンス
//
// @returns
// isExist true 存在, false 存在しない
func (mail *Email) Exist(ctx context.Context, db *repository.Queryer) (bool, error) {
	_, err := mail.repository.FindUserByEmail(ctx, *db, &mail.value)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
