package testutil

import (
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// テストのアサーション関数
// @params
// t テストユーティリティ
// got 実際のレスポンス
// status  wantレスポンスステータスコード
// body wantレスポンスボディ
func AssertResponse(t *testing.T, got *http.Response, status int, body []byte) {
	t.Helper()
	t.Cleanup(func() { _ = got.Body.Close() })
	gb, err := io.ReadAll(got.Body)
	assert.NoError(t, err)

	assert.Equalf(t, status, got.StatusCode, "want http.status %d, but got %d", status, got.StatusCode)

	if len(gb) == 0 && len(body) == 0 {
		// 期待としても実体としてもレスポンスボディがないので
		// AssertJSONを呼ぶ必要はない。
		return
	}
	assert.JSONEq(t, string(body), string(gb))
}

// JSONファイルの読み込み
func LoadFile(t *testing.T, path string) []byte {
	t.Helper()

	bt, err := os.ReadFile(path)
	assert.NoErrorf(t, err, "cannot read from %q", path)
	return bt
}
