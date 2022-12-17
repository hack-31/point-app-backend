package user

import (
	"fmt"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
)

// パスワードオブジェクト
type Password struct {
	value string
}

// パスワードオブジェクト作成
// コンストラクタ
//
// @params pwd パスワード
//
// @return パスワードオブジェクト
func NewPasswrod(pwd string) (*Password, error) {
	if 50 < utf8.RuneCountInString(pwd) {
		return nil, fmt.Errorf("cannot use password over 51 char")
	}
	return &Password{value: pwd}, nil
}

func (pwd *Password) CreateHash() (string, error) {
	pw, err := bcrypt.GenerateFromPassword([]byte(pwd.value), bcrypt.DefaultCost)
	return string(pw), err
}
