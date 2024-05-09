package model

import (
	"math/rand"
	"time"
	"unicode/utf8"

	"github.com/cockroachdb/errors"
	"github.com/hack-31/point-app-backend/constant"
	"golang.org/x/crypto/bcrypt"
)

// パスワードオブジェクト
type Password struct {
	value string
}

// パスワードオブジェクト作成
// ハッシュ化されてない値を扱う
// コンストラクタ
//
// @params pwd パスワード
//
// @return パスワードオブジェクト
func NewPassword(pwd string) (*Password, error) {
	if 50 < utf8.RuneCountInString(pwd) {
		return nil, errors.Wrap(errors.New("cannot use password over 51 char"), "NewPassword")
	}
	return &Password{value: pwd}, nil
}

// ハッシュ化されたパスワードと一致するか
// @params
// hashPwd ハッシュ化されたパスワード
func (pwd *Password) IsMatch(hashPwd string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashPwd), []byte(pwd.value))
	if err != nil {
		return false, err
	}
	return true, err
}

func (pwd *Password) CreateHash() (string, error) {
	pw, err := bcrypt.GenerateFromPassword([]byte(pwd.value), bcrypt.DefaultCost)
	return string(pw), err
}

// ラムダム文字列のパスワードを作成
func (pwd *Password) CreateRandomPassword() *Password {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	value := make([]byte, constant.RandomPasswordLength)
	for i := range value {
		r := r.Int63() % int64(len(letters))
		value[i] = letters[int(r)]
	}
	return &Password{value: string(value)}
}

// 文字列
// @return
// パスワード
func (pwd *Password) String() string {
	return string(pwd.value)
}
