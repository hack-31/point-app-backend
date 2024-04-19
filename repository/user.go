package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/hack-31/point-app-backend/repository/entity"
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
// columns カラム（無指定の場合は全て）
//
// @returns
// entity.User ユーザ情報
func (r *Repository) FindUserByEmail(ctx context.Context, db Queryer, email string, columns ...string) (entity.User, error) {
	formattedColumns := "*"
	if len(columns) > 0 {
		// 以下のような文字列にする
		// id, name, email, created_at, updated_at
		formattedColumns = strings.Join(columns, ", ")
		formattedColumns = strings.TrimSuffix(formattedColumns, ", ")
	}

	sql :=
		`SELECT ` + formattedColumns + `
		FROM users 
		WHERE email = ? 
		LIMIT 1;`

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

// ユーザIDでユーザが存在するか検索する
// @params
// ctx context
// db dbインスタンス
// ID ユーザID
//
// @returns
// entity.User ユーザ情報
func (r *Repository) GetUserByID(ctx context.Context, db Queryer, ID entity.UserID) (entity.User, error) {
	sql := `
		SELECT *
		FROM users
		WHERE id = ?
		LIMIT 1
	`

	var user entity.User
	if err := db.GetContext(ctx, &user, sql, ID); err != nil {
		// 見つけられない時(その他のエラーも含む)
		// 見つけられない時のエラーは利用側で
		// errors.Is(err, sql.ErrNoRows)
		// で判断する
		return user, err
	}
	return user, nil
}

// 削除
func (r *Repository) DeleteUserByID(ctx context.Context, db Execer, ID entity.UserID) (int64, error) {
	sql := "DELETE FROM `users`" +
		"WHERE `id` = ?"
	res, err := db.ExecContext(ctx, sql, ID)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// パスワードを上書き
// @params
// ctx context
// db dbインスタンス
// email email
// pass password
//
// @returns
// error
func (r *Repository) UpdatePassword(ctx context.Context, db Execer, email, pass *string) error {
	sql := `UPDATE users SET password = ? WHERE email = ?`

	_, err := db.ExecContext(ctx, sql, pass, email)
	if err != nil {
		return err
	}
	return err
}

// メールアドレスを上書きする
// @params
// ctx context
// db dbインスタンス
// userID 更新するユーザーID
// newEmail 上書きするメールアドレス
//
// @returns
// error
func (r *Repository) UpdateEmail(ctx context.Context, db Execer, userID entity.UserID, newEmail string) error {
	sql := `UPDATE users SET email = ? WHERE id = ?`

	_, err := db.ExecContext(ctx, sql, newEmail, userID)
	if err != nil {
		return fmt.Errorf("failed to update DB: %w", err)
	}
	return nil
}

// アカウント情報を上書きする
// @params
// ctx context
// db dbインスタンス
// email email
// familyName familyName
// familyNameKana familyNameKana
// firstName firstName
// firstNameKana firstNameKana
//
// @returns
// error
func (r *Repository) UpdateAccount(ctx context.Context, db Execer, email, familyName, familyNameKana, firstName, firstNameKana *string) error {
	sql := `UPDATE users
						SET family_name = ?,
								family_name_kana = ?,
								first_name = ?,
								first_name_kana = ?
						WHERE email = ?`
	_, err := db.ExecContext(
		ctx,
		sql,
		familyName,
		familyNameKana,
		firstName,
		firstNameKana,
		email)
	if err != nil {
		return fmt.Errorf("failed to update DB: %w", err)
	}
	return nil
}

// ユーザ一覧
//
// @params
// ctx context
// db db
// columns 取得カラム（無指定の場合は全て）
//
// @returns
// Users ユーザ一覧
func (r *Repository) GetAll(ctx context.Context, db Queryer, columns ...string) (entity.Users, error) {
	formattedColumns := "*"
	if len(columns) > 0 {
		// 以下のような文字列にする
		// id, name, email, created_at, updated_at
		formattedColumns = strings.Join(columns, ", ")
		formattedColumns = strings.TrimSuffix(formattedColumns, ", ")
	}
	sql :=
		`SELECT ` + formattedColumns + `
		 FROM users;`
	var users entity.Users
	if err := db.SelectContext(ctx, &users, sql); err != nil {
		return users, err
	}
	return users, nil
}
