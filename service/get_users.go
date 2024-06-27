package service

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/hack-31/point-app-backend/domain"
	"github.com/hack-31/point-app-backend/domain/model"
	"github.com/hack-31/point-app-backend/repository"
	"github.com/jmoiron/sqlx"
)

type GetUsers struct {
	DB              repository.Queryer
	UserRepo        domain.UserRepo
	TransactionRepo domain.TransactionRepo
	TokenGenerator  domain.TokenGenerator
}

func NewGetUsers(db *sqlx.DB, repo *repository.Repository, jwter domain.TokenGenerator) *GetUsers {
	return &GetUsers{
		DB:              db,
		UserRepo:        repo,
		TransactionRepo: repo,
		TokenGenerator:  jwter,
	}
}

type GetUsersResponse struct {
	Users []struct {
		ID               model.UserID
		FirstName        string
		FirstNameKana    string
		FamilyName       string
		FamilyNameKana   string
		Email            string
		AcquisitionPoint int
	}
}

// ユーザ一覧取得サービス
//
// @params ctx コンテキスト
//
// @return
// ユーザ一覧
func (r *GetUsers) GetUsers(ctx context.Context) (GetUsersResponse, error) {
	// ユーザ一覧を取得する
	users, err := r.UserRepo.GetAll(ctx, r.DB)
	if err != nil {
		return GetUsersResponse{}, errors.Wrap(err, "failed to get users in GetUsersService.GetUsers")
	}

	// ユーザIDsを取得する
	userIDs := make([]model.UserID, 0, len(users))
	for _, user := range users {
		userIDs = append(userIDs, model.UserID(user.ID))
	}

	// 取得ポイントを取得する
	points, err := r.TransactionRepo.GetAquistionPoint(ctx, r.DB, userIDs)
	if err != nil {
		return GetUsersResponse{}, errors.Wrap(err, "failed to get points in GetUsersService.GetUsers")
	}

	res := make([]struct {
		ID               model.UserID
		FirstName        string
		FirstNameKana    string
		FamilyName       string
		FamilyNameKana   string
		Email            string
		AcquisitionPoint int
	}, 0, len(users))

	// ユーザに取得ポイントを設定する
	for _, v := range users {
		res = append(res, struct {
			ID               model.UserID
			FirstName        string
			FirstNameKana    string
			FamilyName       string
			FamilyNameKana   string
			Email            string
			AcquisitionPoint int
		}{
			ID:               model.UserID(v.ID),
			FirstName:        v.FirstName,
			FirstNameKana:    v.FirstNameKana,
			FamilyName:       v.FamilyName,
			FamilyNameKana:   v.FamilyNameKana,
			Email:            v.Email,
			AcquisitionPoint: points[model.UserID(v.ID)],
		})
	}
	return GetUsersResponse{
		Users: res,
	}, nil
}
