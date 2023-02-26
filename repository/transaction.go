package repository

import (
	"context"
	"database/sql"
	"fmt"
)

type AppConnection struct {
	// DBインスタンス
	db Beginner
	// トランザクションで利用するインスタンス
	Tx *sql.Tx
}

// トランザクション
func NewAppConnection(db Beginner) *AppConnection {
	return &AppConnection{db: db}
}

// トラザクション開始
func (ac *AppConnection) Begin(ctx context.Context) error {
	tx, err := ac.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("cannet connect transaction: %w", err)
	}
	ac.Tx = tx
	return nil
}

// コミット
// トランザクションの最後に実行
func (ac *AppConnection) Commit() error {
	return ac.Tx.Commit()
}

// ロールバック
// トラザクションを開いてから、エラーが起きた時に実行する
func (ac *AppConnection) Rollback() error {
	return ac.Tx.Rollback()
}
