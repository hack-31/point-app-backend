package service

import (
	"context"
	"fmt"

	"github.com/hack-31/point-app-backend/constant"
	"github.com/hack-31/point-app-backend/domain/user"
	"github.com/hack-31/point-app-backend/entity"
	"github.com/hack-31/point-app-backend/repository"
)

type RegisterUser struct {
	DB    repository.Execer
	Cache *repository.KVS
	Repo  UserRegister
}

// ユーザ登録サービス
//
// @params ctx コンテキスト, temporaryUserId 一時保存ユーザid
//
// @return ユーザ情報
func (r *RegisterUser) RegisterUser(ctx context.Context, temporaryUserId string) (*entity.User, error) {
	// 一時ユーザ情報を復元
	u, err := r.Cache.Load(ctx, temporaryUserId)
	if err != nil {
		return nil, fmt.Errorf("cannot load user in redis: %w", err)
	}

	// 復元が成功したら一時ユーザ情報除削
	if err := r.Cache.Delete(ctx, temporaryUserId); err != nil {
		return nil, fmt.Errorf("cannot delete in redis: %w", err)
	}

	// 復元したユーザ情報を解析
	temporyUser := user.NewTemporaryUserString(u)
	firstName, firstNameKana, familyName, familyNameKana, email, hashPass := temporyUser.Split()

	// DBに保存
	user := &entity.User{
		FirstName:      firstName,
		FirstNameKana:  firstNameKana,
		FamilyName:     familyName,
		FamilyNameKana: familyNameKana,
		Email:          email,
		Password:       hashPass,
		SendingPoint:   constant.DefaultSendingPoint,
	}
	if err := r.Repo.RegisterUser(ctx, r.DB, user); err != nil {
		return nil, fmt.Errorf("failed to register: %w", err)
	}

	// TODO: セッションキー(token)作成

	return user, nil
}
