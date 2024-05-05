package myerror

import (
	"github.com/cockroachdb/errors"
)

var (
	// 認証エラー
	ErrNotExistEmail = errors.New("メールアドレスが存在しません。")
	// 登録済みのメールアドレスは登録できません
	ErrAlreadyEntry = errors.New("登録済みのメールアドレスは登録できません。")
	// 確認コードまたは、セッションキーが無効です
	ErrNotFoundSession = errors.New("確認コードまたは、セッションキーが無効です。")
	// メールアドレスまたは、パスワードが異なるエラー
	ErrNotMatchLogInfo = errors.New("メールアドレスまたは、パスワードが異なります。")
	//　ユーザが存在しないエラー
	ErrNotUser = errors.New("ユーザが存在しません。")
	//　送付可能ポイントが不足しているエラー
	ErrHasNotSendablePoint = errors.New("送付可能ポイントが不足しています。")
	// パスワードが異なるエラー
	ErrDifferentPassword = errors.New("パスワードが異なります。")
	// データが存在しないエラー
	ErrNotFound = errors.New("データが存在しません。")
	// データベースエラー
	ErrDBException = errors.New("データベースエラー")
	// キャッシュエラー
	ErrCacheException = errors.New("キャッシュエラー")
)
