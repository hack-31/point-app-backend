package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type AppConnection struct {
	// DBインスタンス
	db Beginner
	// トランザクションで利用するインスタンス
	tx *sqlx.Tx
}

// トランザクション
func NewAppConnection(db Beginner) *AppConnection {
	return &AppConnection{db: db}
}

// トラザクション開始
func (ac *AppConnection) Begin(ctx context.Context) error {
	tx, err := ac.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("cannet connect transaction: %w", err)
	}
	ac.tx = tx
	return nil
}

// コミット
// トランザクションの最後に実行
func (ac *AppConnection) Commit() error {
	if err := ac.tx.Commit(); err != nil {
		return fmt.Errorf("cannot commit: %w ", err)
	}
	return nil
}

// ロールバック
// トラザクションを開いてから、エラーが起きた時に実行する
func (ac *AppConnection) Rollback() error {
	if err := ac.tx.Rollback(); err != nil {
		return fmt.Errorf("cannot rollback: %w", err)
	}
	return nil
}

// トランザクション用DBインスタンス
func (ac *AppConnection) DB() *sqlx.Tx {
	return ac.tx
}
