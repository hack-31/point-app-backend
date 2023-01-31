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
// @params
// ctx コンテキスト
// db dbの値(インスタンス)
// u ユーザエンティティ
func (r *Repository) RegisterUser(ctx context.Context, db Execer, u *entity.User) error {
	u.CreatedAt = r.Clocker.Now()
	u.UpdateAt = r.Clocker.Now()

	sql := `INSERT INTO users (
			first_name, first_name_kana, family_name, family_name_kana, email, password, sending_point, created_at, update_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := db.ExecContext(ctx, sql, u.FirstName, u.FirstNameKana, u.FamilyName, u.FamilyNameKana, u.Email, u.Password, u.SendingPoint, u.CreatedAt, u.UpdateAt)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == ErrCodeMySQLDuplicateEntry {
			return fmt.Errorf("cannot create same email user: %w", ErrAlreadyEntry)
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

// メールでユーザが存在するか検索する
// @params
// ctx context
// db dbインスタンス
// email email
//
// @returns
// entity.User ユーザ情報
func (r *Repository) FindUserByEmail(ctx context.Context, db Queryer, email *string) (entity.User, error) {
	sql := `
		SELECT 
			u.id,
			u.first_name, 
			u.first_name_kana, 
			u.family_name, 
			u.family_name_kana, 
			u.email,
			u.password,
			u.created_at,
			u.update_at,
			u.sending_point,
			SUM(IFNULL(t.transaction_point, 0)) AS acquisition_point 
		from users AS u
		LEFT JOIN transactions AS t
		ON u.id = t.receiving_user_id
		WHERE u.email = ?
		GROUP BY u.id
		LIMIT 1`

	var user entity.User

	if err := db.GetContext(ctx, &user, sql, email); err != nil {
		// 見つけられない時(その他のエラーも含む)
		// 見つけられない時のエラーは利用側で
		// errors.Is(err, sql.ErrNoRows)
		// で判断する
		return user, err
	}
	return user, nil
}

// ユーザ一覧
//
// @params
// ctx context
// db db
//
// @returns
// Users ユーザ一覧
func (r *Repository) FindUsers(ctx context.Context, db Queryer) (entity.Users, error) {
	sql := `
		SELECT 
			u.id, 
			u.first_name, 
			u.first_name_kana, 
			u.family_name, 
			u.family_name_kana, 
			u.email, 
			SUM(IFNULL(t.transaction_point, 0)) AS acquisition_point 
		from users AS u
		LEFT JOIN transactions AS t
		ON u.id = t.receiving_user_id
		GROUP BY u.id`

	var users entity.Users
	if err := db.SelectContext(ctx, &users, sql); err != nil {
		return users, err
	}
	return users, nil
}
