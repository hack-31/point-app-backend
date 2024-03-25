package repository

import (
	"context"

	"github.com/hack-31/point-app-backend/domain/model"
	"github.com/jmoiron/sqlx"
)

// GetAquistionPoint は、指定ユーザーの取得ポイントを取得する
func (r *Repository) GetAquistionPoint(ctx context.Context, db Queryer, userIDs []model.UserID) (map[model.UserID]int, error) {
	query := `
		SELECT
			receiving_user_id,
			SUM(transaction_point) as acquisition_point
		FROM
			transactions
		WHERE
			receiving_user_id IN (?)
		GROUP BY
			receiving_user_id;
		`

	// クエリ作成
	query, params, err := sqlx.In(query, userIDs)
	if err != nil {
		return nil, err
	}

	// 取得
	var users []struct {
		ID               model.UserID `db:"receiving_user_id"`
		AcquisitionPoint int          `db:"acquisition_point"`
	}
	if err := db.SelectContext(ctx, &users, query, params...); err != nil {
		return nil, err
	}

	// データ整形
	acquistionPoints := make(map[model.UserID]int, len(users))
	for _, v := range users {
		acquistionPoints[v.ID] = v.AcquisitionPoint
	}
	return acquistionPoints, nil
}
