package repository

import (
	"errors"
)

const (
	// ErrCodeMySQLDuplicateEntry はMySQL系のDUPLICATEエラーコード
	// https://dev.mysql.com/doc/mysql-errors/8.0/en/server-error-reference.html
	// Error number: 1062; Symbol: ER_DUP_ENTRY; SQLSTATE: 23000
	ErrCodeMySQLDuplicateEntry  = 1062
	ErrCodeMySQLNoReferencedRow = 1452
)

var (
	ErrNotExistEmail       = errors.New("メールアドレスが存在しません。")
	ErrAlreadyEntry        = errors.New("登録済みのメールアドレスは登録できません。")
	ErrNotFoundSession     = errors.New("確認コードまたは、セッションキーが無効です。")
	ErrNotMatchLogInfo     = errors.New("メールアドレスまたは、パスワードが異なります。")
	ErrNotUser             = errors.New("ユーザが存在しません。")
	ErrHasNotSendablePoint = errors.New("送付可能ポイントが不足しています。")
)
