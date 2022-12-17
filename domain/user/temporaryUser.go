package user

import (
	"fmt"
	"regexp"
)

type TemporaryUserString struct {
	value string
}

func NewTemporaryUserString(temporaryUserString string) *TemporaryUserString {
	return &TemporaryUserString{value: temporaryUserString}
}

// ユーザ情報を開業で区切り、1つの文字列に結合する
//
// @params
// ctx コンテキスト
// firstName 名前
// firstNameKana 名前カナ
// familyName 名字
// familyNameKana 名字カナ
// email メールアドレス
// password パスワード
//
// @returns
// 連結したユーザ
func (tus *TemporaryUserString) Join(firstName, firstNameKana, familyName, familyNameKana, email, password string) string {
	tus.value = fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s", firstName, firstNameKana, familyName, familyNameKana, email, password)
	return tus.value
}

// 開業で区切られた1つの文字列になっているユーザ情報を分解する
//
// @returns
// ctx コンテキスト
// firstName 名前
// firstNameKana 名前カナ
// familyName 名字
// familyNameKana 名字カナ
// email メールアドレス
// hashPass パスワード
func (tus *TemporaryUserString) Split() (firstName, firstNameKana, familyName, familyNameKana, email, hashPass string) {
	reg := "\r\n|\n"
	arr := regexp.MustCompile(reg).Split(tus.value, -1)
	return arr[0], arr[1], arr[2], arr[3], arr[4], arr[5]
}
