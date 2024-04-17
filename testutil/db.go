package testutil

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

// テストにおいてデータベース接続する
// 実際にデーターベース接続する
func OpenDBForTest(t *testing.T) *sqlx.DB {
	t.Helper()

	host := "db"
	if _, defined := os.LookupEnv("CI"); defined {
		host = "127.0.0.1"
	}

	db, err := sql.Open(
		"mysql",
		fmt.Sprintf("admin:password@tcp(%s:3306)/point_app_test?parseTime=true&loc=Asia%%2FTokyo", host),
	)
	assert.NoError(t, err, "cannot open db")

	t.Cleanup(
		func() { _ = db.Close() },
	)

	xdb := sqlx.NewDb(db, "mysql")

	// AUTO_INCREMENTをリセットする
	_, err = xdb.Exec("ALTER TABLE `users` AUTO_INCREMENT = 1;")
	assert.NoError(t, err)

	return xdb
}

// NewDBForTest はテスト用のDBを作成する
//
// 実際にはDBには接続しない
func NewDBForMock(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock) {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		assert.Fail(t, err.Error())
	}

	xdb := sqlx.NewDb(db, "sqlmock")
	t.Cleanup(
		func() { _ = xdb.Close() },
	)
	if err := xdb.Ping(); err != nil {
		assert.Error(t, err)
	}

	return xdb, mock
}

// NewTxForMock はモック用のトランザクションを作成する
//
// 実際にはDBには接続しない
//
// begin, commitは正常系を返すようにモックする
// 異常系をテストしたい場合は、NewDBForMockを使って自分でモックを作成する
func NewTxForMock(t *testing.T, ctx context.Context) *sqlx.Tx {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		assert.Fail(t, err.Error())
	}

	xdb := sqlx.NewDb(db, "sqlmock")
	t.Cleanup(
		func() { _ = xdb.Close() },
	)
	if err := xdb.Ping(); err != nil {
		assert.Error(t, err)
	}

	// トランザクションの開始、コミットをモックする
	// 正常系を返す
	mock.ExpectBegin()
	mock.ExpectCommit()

	mockTx, err := xdb.BeginTxx(ctx, nil)
	if err != nil {
		assert.Error(t, err)
	}
	return mockTx
}
