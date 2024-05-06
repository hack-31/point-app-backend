package service

import (
	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
	"github.com/hack-31/point-app-backend/domain"
	"github.com/hack-31/point-app-backend/repository"
	"github.com/hack-31/point-app-backend/repository/entity"
	"github.com/hack-31/point-app-backend/utils"
	"github.com/jmoiron/sqlx"
)

type GetAccount struct {
	DB              repository.Queryer
	Repo            domain.UserRepo
	TransactionRepo domain.TransactionRepo
}

func NewGetAccount(db *sqlx.DB, repo *repository.Repository) *GetAccount {
	return &GetAccount{
		DB:              db,
		Repo:            repo,
		TransactionRepo: repo,
	}
}

type GetAccountResponse struct {
	ID               entity.UserID
	FirstName        string
	FirstNameKana    string
	FamilyName       string
	FamilyNameKana   string
	Email            string
	AcquisitionPoint int
	SendablePoint    int
}

// ユーザ一覧取得サービス
//
// @params ctx コンテキスト
//
// @return
// ユーザ一覧
func (ga *GetAccount) GetAccount(ctx *gin.Context) (GetAccountResponse, error) {
	mail := utils.GetEmail(ctx)

	// Emailよりユーザ情報を取得する
	user, err := ga.Repo.FindUserByEmail(ctx, ga.DB, mail)
	if err != nil {
		return GetAccountResponse{}, errors.Wrap(err, "failed to get user by email")
	}

	// 取得ポイントを取得する
	userID := utils.GetUserID(ctx)
	point, err := ga.TransactionRepo.GetAquistionPoint(ctx, ga.DB, []entity.UserID{userID})
	if err != nil {
		return GetAccountResponse{}, errors.Wrap(err, "failed to get acquisition point")
	}

	return GetAccountResponse{
		ID:               user.ID,
		FirstName:        user.FirstName,
		FirstNameKana:    user.FirstNameKana,
		FamilyName:       user.FamilyName,
		FamilyNameKana:   user.FamilyNameKana,
		Email:            user.Email,
		SendablePoint:    user.SendingPoint,
		AcquisitionPoint: point[userID],
	}, nil
}
