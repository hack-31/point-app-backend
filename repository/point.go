package repository

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/go-sql-driver/mysql"
	"github.com/hack-31/point-app-backend/myerror"
	"github.com/hack-31/point-app-backend/repository/entity"
)

// ポイントの取引履歴の保存
//
// @params
// ctx コンテキスト
// db dbの値(インスタンス)
// fromUserID 送信元ユーザ
// toUserId 送信先ユーザ
// sendPoint 送付ポイント
func (r *Repository) RegisterPointTransaction(ctx context.Context, db Execer, fromUserID, toUserId entity.UserID, sendPoint int) error {
	sql := `INSERT INTO transactions (
		sending_user_id, receiving_user_id, transaction_point, transaction_at
	) VALUES (?, ?, ?, ?)`

	if _, err := db.ExecContext(ctx, sql, fromUserID, toUserId, sendPoint, r.Clocker.Now()); err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == ErrCodeMySQLNoReferencedRow {
			return errors.Join(myerror.ErrNotUser, err)
		}
		return err
	}

	return nil
}

// ポイント可能額の更新
//
// @params
// ctx コンテキスト
// db db_instanc
// fromUserID 送付者ユーザのID
// point ポイント残高
func (r *Repository) UpdateSendablePoint(ctx context.Context, db Execer, fromUserID entity.UserID, point int) error {
	sql := `
	  UPDATE users
		SET sending_point = ?
		WHERE id = ?
	`
	if _, err := db.ExecContext(ctx, sql, point, fromUserID); err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == ErrCodeMySQLNoReferencedRow {
			return errors.Join(myerror.ErrNotUser, err)
		}
		return err
	}
	return nil
}

// UpdateAllSendablePoint は、全ユーザの送付可能ポイントを更新する
func (r *Repository) UpdateAllSendablePoint(ctx context.Context, db Execer, point int) error {
	sql := `
	  UPDATE users SET sending_point = ?;
	`
	if _, err := db.ExecContext(ctx, sql, point); err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == ErrCodeMySQLNoReferencedRow {
			return errors.Join(myerror.ErrNotUser, err)
		}
		return err
	}
	return nil
}
