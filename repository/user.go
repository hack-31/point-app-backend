package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/hack-31/point-app-backend/entity"
)

// ユーザ情報永続化
//
// @params ctx コンテキスト
//
// @paramas db dbの値(インスタンス)
//
// @params u ユーザエンティティ
func (r *Repository) RegisterUser(ctx context.Context, db Execer, u *entity.User) error {
	u.CreatedAt = r.Clocker.Now()
	u.ModifiedAt = r.Clocker.Now()
	sql := `INSERT INTO users (
			name, email, password, role_id, created_at, modified_at
		) VALUES (?, ?, ?, ?, ?, ?)`
	result, err := db.ExecContext(ctx, sql, u.Name, u.Email, u.Password, u.Role, u.CreatedAt, u.ModifiedAt)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == ErrCodeMySQLDuplicateEntry {
			return fmt.Errorf("cannot create same name user: %w", ErrAlreadyEntry)
		}
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	u.ID = entity.UserID(id)
	return nil
}
