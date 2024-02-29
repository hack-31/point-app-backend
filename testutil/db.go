package testutil

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

// テストにおいてデータベース接続する
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
	xdb.Exec("ALTER TABLE `users` AUTO_INCREMENT = 1;")

	return xdb
}
