package testdata

import (
	"context"
	"testing"

	"github.com/hack-31/point-app-backend/repository/entity"
	"github.com/hack-31/point-app-backend/utils/clock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

// 取引テーブルにデータを挿入する
func Transactions(t *testing.T, ctx context.Context, con *sqlx.Tx, setter func(users entity.Transactions)) entity.Transactions {
	t.Helper()
	if _, err := con.ExecContext(ctx, `DELETE FROM transactions;`); err != nil {
		t.Logf("%v", err)
	}

	c := clock.FixedClocker{
		IsAsia: true,
	}

	// デフォルト値
	transactions := entity.Transactions{
		{
			ReceivingUserID:  1,
			SendingUserID:    2,
			TransactionPoint: 1000,
			TransactionAt:    c.Now(),
		},
		{
			ReceivingUserID:  1,
			SendingUserID:    2,
			TransactionPoint: 200,
			TransactionAt:    c.Now(),
		},
		{
			ReceivingUserID:  2,
			SendingUserID:    1,
			TransactionPoint: 500,
			TransactionAt:    c.Now(),
		},
		{
			ReceivingUserID:  3,
			SendingUserID:    2,
			TransactionPoint: 300,
			TransactionAt:    c.Now(),
		},
	}

	// データの上書き
	setter(transactions)

	// データ道入
	result, err := con.NamedExecContext(ctx, `
		INSERT INTO transactions
			(receiving_user_id, sending_user_id, transaction_point, transaction_at)
		VALUES
			(:receiving_user_id, :sending_user_id, :transaction_point, :transaction_at);`,
		transactions,
	)
	assert.NoError(t, err)

	// 発行したIDを取得
	id, err := result.LastInsertId()
	assert.NoError(t, err)

	// IDの挿入
	transactions[0].ID = entity.TransactionID(id)
	transactions[1].ID = entity.TransactionID(id + 1)
	transactions[2].ID = entity.TransactionID(id + 2)
	transactions[3].ID = entity.TransactionID(id + 2)

	return transactions
}
