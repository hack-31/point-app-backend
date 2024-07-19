package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/hack-31/point-app-backend/domain"
	"github.com/hack-31/point-app-backend/domain/model"
	"github.com/hack-31/point-app-backend/repository"
	"github.com/hack-31/point-app-backend/repository/entities"
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

type GetUsersRequest struct {
	Size       int
	NextCursor string
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
	NextCursor string
}

// ユーザ一覧取得サービス
//
// @params ctx コンテキスト
//
// @return
// ユーザ一覧
func (r *GetUsers) GetUsers(ctx context.Context, input GetUsersRequest) (GetUsersResponse, error) {
	// ユーザ一覧を取得する
	type cursor struct {
		UserID    model.UserID `json:"user_id"`
		Point     int          `json:"point"`
		CreatedAt time.Time    `json:"created_at"`
	}
	var users []*entities.User
	var c cursor
	if input.NextCursor != "" {
		data, err := base64.URLEncoding.DecodeString(input.NextCursor)
		if err != nil {
			return GetUsersResponse{}, errors.Wrap(err, "failed to decode nextCursor in GetUsersService.GetUsers")
		}
		if err := json.Unmarshal(data, &c); err != nil {
			return GetUsersResponse{}, errors.Wrap(err, "failed to unmarshal nextCursor in GetUsersService.GetUsers")
		}
		users, err = r.UserRepo.GetAllWithCursor(ctx, r.DB, repository.GetAllWithCursorParam{
			Size:            input.Size,
			CursorPoint:     c.Point,
			CursorUserID:    c.UserID,
			CursorCreatedAt: c.CreatedAt,
		})
		if err != nil {
			return GetUsersResponse{}, errors.Wrap(err, "failed to get users in GetUsersService.GetUsers")
		}
	}
	if input.NextCursor == "" {
		var err error
		if input.Size == 0 {
			input.Size = 10
		}
		users, err = r.UserRepo.GetUsers(ctx, r.DB, repository.GetUsersParam{
			Size: input.Size,
		})
		if err != nil {
			return GetUsersResponse{}, errors.Wrap(err, "failed to get users in GetUsersService.GetUsers")
		}
	}
	if len(users) == 0 {
		return GetUsersResponse{
			Users: []struct {
				ID               model.UserID
				FirstName        string
				FirstNameKana    string
				FamilyName       string
				FamilyNameKana   string
				Email            string
				AcquisitionPoint int
			}{},
			NextCursor: "",
		}, nil
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
	var nextCursorStr string
	if len(users) == input.Size {
		data, err := json.Marshal(cursor{
			UserID:    model.UserID(users[len(users)-1].ID),
			Point:     points[model.UserID(users[len(users)-1].ID)],
			CreatedAt: users[len(users)-1].CreatedAt,
		})
		if err != nil {
			return GetUsersResponse{}, errors.Wrap(err, "failed to marshal nextCursor in GetUsersService.GetUsers")
		}
		nextCursorStr = base64.URLEncoding.EncodeToString(data)
	}

	return GetUsersResponse{
		Users:      res,
		NextCursor: nextCursorStr,
	}, nil
}
