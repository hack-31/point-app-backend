package repository

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"
	"github.com/hack-31/point-app-backend/domain/model"
	"github.com/hack-31/point-app-backend/testutil"
	"github.com/hack-31/point-app-backend/utils/clock"
	"github.com/jmoiron/sqlx"
)

func prepareUsers(ctx context.Context, t *testing.T, con Execer) model.Users {
	t.Helper()
	// 外部キー制約を一時的ににする
	if _, err := con.ExecContext(ctx, `SET FOREIGN_KEY_CHECKS=0;`); err != nil {
		t.Logf("%v", err)
	}
	// 一旦データを全て削除
	if _, err := con.ExecContext(ctx, `DELETE FROM transactions;`); err != nil {
		t.Logf("%v", err)
	}
	if _, err := con.ExecContext(ctx, `DELETE FROM users;`); err != nil {
		t.Logf("%v", err)
	}

	c := clock.FixedClocker{}
	// 全ユーザ取得メソッドを呼び出した時の期待データ
	wants := model.Users{
		{
			FirstName:      "太郎",
			FirstNameKana:  "たろう",
			FamilyName:     "山田",
			FamilyNameKana: "やまだ",
			Email:          "yamada.taro@example.com",
		},
		{
			FirstName:      "葵",
			FirstNameKana:  "あおい",
			FamilyName:     "斉藤",
			FamilyNameKana: "さいとう",
			Email:          "saito.aoi@example.com",
		},
		{
			FirstName:      "拓也",
			FirstNameKana:  "たくや",
			FamilyName:     "木村",
			FamilyNameKana: "きむら",
			Email:          "kimura.takuya@example.com",
		},
	}
	// 挿入データ
	users := model.Users{
		{
			FirstName:      wants[0].FirstName,
			FirstNameKana:  wants[0].FirstNameKana,
			FamilyName:     wants[0].FamilyName,
			FamilyNameKana: wants[0].FamilyNameKana,
			Email:          wants[0].Email,
			UpdateAt:       c.Now(),
			CreatedAt:      c.Now(),
		},
		{
			FirstName:      wants[1].FirstName,
			FirstNameKana:  wants[1].FirstNameKana,
			FamilyName:     wants[1].FamilyName,
			FamilyNameKana: wants[1].FamilyNameKana,
			Email:          wants[1].Email,
			UpdateAt:       c.Now(),
			CreatedAt:      c.Now(),
		},
		{
			FirstName:      wants[2].FirstName,
			FirstNameKana:  wants[2].FirstNameKana,
			FamilyName:     wants[2].FamilyName,
			FamilyNameKana: wants[2].FamilyNameKana,
			Email:          wants[2].Email,
			UpdateAt:       c.Now(),
			CreatedAt:      c.Now(),
		},
	}
	result, err := con.ExecContext(ctx,
		`INSERT INTO users (first_name, first_name_kana, family_name, family_name_kana, email, password, sending_point, created_at, update_at)
			VALUES
			    (?, ?, ?, ?, ?, ?, ?, ?, ?),
			    (?, ?, ?, ?, ?, ?, ?, ?, ?),
			    (?, ?, ?, ?, ?, ?, ?, ?, ?);`,
		users[0].FirstName, users[0].FirstNameKana, users[0].FamilyName, users[0].FamilyNameKana, users[0].Email, users[0].Password, users[0].SendingPoint, users[0].CreatedAt, users[0].UpdateAt,
		users[1].FirstName, users[1].FirstNameKana, users[1].FamilyName, users[1].FamilyNameKana, users[1].Email, users[1].Password, users[1].SendingPoint, users[1].CreatedAt, users[1].UpdateAt,
		users[2].FirstName, users[2].FirstNameKana, users[2].FamilyName, users[2].FamilyNameKana, users[2].Email, users[2].Password, users[2].SendingPoint, users[2].CreatedAt, users[2].UpdateAt,
	)
	if err != nil {
		t.Fatal(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}
	wants[0].ID = model.UserID(id)
	wants[1].ID = model.UserID(id + 1)
	wants[2].ID = model.UserID(id + 2)
	return wants
}

func TestRepository_FindUsers(t *testing.T) {
	t.Parallel()
	t.Run("ユーザ一覧を取得する", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		// model.Usersを作成する他のテストケースと混ざるとテストがフェイルする。
		// そのため、トランザクションをはることでこのテストケースの中だけのテーブル状態にする。
		tx, err := testutil.OpenDBForTest(t).BeginTxx(ctx, nil)
		// このテストケースが完了したらもとに戻す
		t.Cleanup(func() {
			_ = tx.Rollback()
			// 外部キー制約を有効にする
			_, _ = tx.ExecContext(ctx, `SET FOREIGN_KEY_CHECKS=1`)
		})
		if err != nil {
			t.Fatal(err)
		}

		wants := prepareUsers(ctx, t, tx)

		// 実行
		r := &Repository{}
		gots, err := r.FindUsers(ctx, tx)

		// アサーション
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if d := cmp.Diff(gots, wants); len(d) != 0 {
			t.Errorf("differs: (-got +want)\n%s", d)
		}
	})
}

func TestRepository_RegisterUser(t *testing.T) {
	t.Parallel()
	t.Run("ユーザを登録し、ユーザIDを取得し、ユーザ情報ポインタにセットする", func(t *testing.T) {
		c := clock.FixedClocker{}
		// 登録データ
		okUser := &model.User{
			FirstName:      "山田",
			FirstNameKana:  "やまだ",
			FamilyName:     "太郎",
			FamilyNameKana: "たろう",
			Password:       "pass",
			Email:          "yamada.taro@example.com",
			CreatedAt:      c.Now(),
			UpdateAt:       c.Now(),
		}

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = db.Close() })

		var wantID int64 = 10
		// モックの定義
		mock.ExpectExec(
			// エスケープが必要
			`INSERT INTO users \(
			first_name, first_name_kana, family_name, family_name_kana, email, password, sending_point, created_at, update_at
		\) VALUES \(\?, \?, \?, \?, \?, \?, \?, \?, \?\)`,
		).WithArgs(
			okUser.FirstName,
			okUser.FirstNameKana,
			okUser.FamilyName,
			okUser.FamilyNameKana,
			okUser.Email,
			okUser.Password,
			okUser.SendingPoint,
			okUser.CreatedAt,
			okUser.UpdateAt,
		).WillReturnResult(sqlmock.NewResult(wantID, 1))

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		xdb := sqlx.NewDb(db, "mysql")
		r := &Repository{Clocker: c}

		if err := r.RegisterUser(ctx, xdb, okUser); err != nil {
			t.Errorf("want no error, but got %v", err)
		}
		if d := cmp.Diff(okUser.ID, model.UserID(wantID)); len(d) != 0 {
			t.Errorf("differs: (-got +want)\n%s", d)
		}
	})
}
