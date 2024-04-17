package usecase

import (
	"context"

	"github.com/hack-31/point-app-backend/domain"
	"github.com/hack-31/point-app-backend/repository"
	"github.com/jmoiron/sqlx"
)

type ResetSendablePoint struct {
	Tx        repository.Beginner
	PointRepo domain.PointRepo
}

func NewResetSendablePoint(db *sqlx.DB, repo *repository.Repository) *ResetSendablePoint {
	return &ResetSendablePoint{
		Tx:        db,
		PointRepo: repo,
	}
}

// ResetPoint は、ポイントをリセットする
func (rsp *ResetSendablePoint) ResetPoint(ctx context.Context, initialSendablePoint int) error {
	tx, err := rsp.Tx.BeginTxx(ctx, nil)
	defer func() { _ = tx.Rollback() }()

	if err != nil {
		return err
	}
	if err := rsp.PointRepo.UpdateAllSendablePoint(ctx, tx, initialSendablePoint); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
