package testutil

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
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
		fmt.Sprintf("admin:password@tcp(%s:3306)/point_app?parseTime=true&loc=Asia%%2FTokyo", host),
	)

	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(
		func() { _ = db.Close() },
	)
	return sqlx.NewDb(db, "mysql")
}
