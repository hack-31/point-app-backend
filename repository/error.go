package repository

import (
	"errors"
)

const (
	// ErrCodeMySQLDuplicateEntry はMySQL系のDUPLICATEエラーコード
	// https://dev.mysql.com/doc/mysql-errors/8.0/en/server-error-reference.html
	// Error number: 1062; Symbol: ER_DUP_ENTRY; SQLSTATE: 23000
	ErrCodeMySQLDuplicateEntry = 1062
)

var (
	ErrAlreadyEntry    = errors.New("登録済みのメールアドレスは登録できません。")
	ErrNotExistEmail   = errors.New("メールアドレスが存在しません。")
	ErrNotFoundSession = errors.New("確認コードまたは、セッションキーが無効です。")
	ErrNotMatchLogInfo = errors.New("メールアドレスまたは、パスワードが異なります。")
)
