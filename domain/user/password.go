package user

import (
	"fmt"
	"math/rand"
	"time"
	"unicode/utf8"

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
func NewPasswrod(pwd string) (*Password, error) {
	if 50 < utf8.RuneCountInString(pwd) {
		return nil, fmt.Errorf("cannot use password over 51 char")
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
	// どの文字列からランダム文字列を生成するか
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	rand.Seed(time.Now().UnixNano())
	value := make([]byte, constant.RandomPasswordLength)
	for i := range value {
		r := rand.Int63() % int64(len(letters))
		value[i] = letters[int(r)]
	}
	return &Password{value: string(value)}
}
