package testdata

import (
	"context"
	"testing"

	"github.com/hack-31/point-app-backend/repository/entity"
	"github.com/hack-31/point-app-backend/utils/clock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

// ユーザーテーブルにデータを挿入する
func Users(t *testing.T, ctx context.Context, con *sqlx.Tx, setter func(users entity.Users)) entity.Users {
	t.Helper()
	// ユーザ一覧を削除する
	if _, err := con.ExecContext(ctx, `DELETE FROM delete_users;`); err != nil {
		t.Logf("%v", err)
	}
	if _, err := con.ExecContext(ctx, `DELETE FROM users;`); err != nil {
		t.Logf("%v", err)
	}

	c := clock.FixedClocker{
		IsAsia: true,
	}

	// ユーザーデータのデフォルト値
	users := entity.Users{
		{
			FirstName:      "太郎",
			FirstNameKana:  "タロウ",
			FamilyName:     "本田",
			FamilyNameKana: "ホンダ",
			Email:          "honda.taro@sample.com",
			Password:       "honda.pass",
			SendingPoint:   0,
			CreatedAt:      c.Now(),
			UpdateAt:       c.Now(),
		},
		{
			FirstName:      "葵",
			FirstNameKana:  "あおい",
			FamilyName:     "斉藤",
			FamilyNameKana: "さいとう",
			Email:          "saito.aoi@example.com",
			Password:       "aoi.pass",
			SendingPoint:   100,
			CreatedAt:      c.Now(),
			UpdateAt:       c.Now(),
		},
		{
			FirstName:      "拓也",
			FirstNameKana:  "たくや",
			FamilyName:     "木村",
			FamilyNameKana: "きむら",
			Email:          "kimura.takuya@example.com",
			Password:       "kimura.pass",
			SendingPoint:   800,
			CreatedAt:      c.Now(),
			UpdateAt:       c.Now(),
		},
	}

	// データの上書き
	setter(users)

	// データ道入
	result, err := con.NamedExecContext(ctx, `
		INSERT INTO users
			(first_name, first_name_kana, family_name, family_name_kana, email, password, sending_point, created_at, update_at)
		VALUES
			(:first_name, :first_name_kana, :family_name, :family_name_kana, :email, :password, :sending_point, :created_at, :update_at);`, users)
	assert.NoError(t, err)

	// 自動発行されたユーザーIDを取得
	id, err := result.LastInsertId()
	assert.NoError(t, err)

	users[0].ID = entity.UserID(id)
	users[1].ID = entity.UserID(id + 1)
	users[2].ID = entity.UserID(id + 2)

	return users
}
