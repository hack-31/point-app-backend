package model

import (
	"math/rand"
	"time"

	"github.com/hack-31/point-app-backend/constant"
)

type ConfirmCode struct {
	value []byte
}

// ラムダム文字列の確認コードを作成
func NewConfirmCode() *ConfirmCode {
	const letters = "0123456789"

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	value := make([]byte, constant.ConfirmCodeLength)
	for i := range value {
		r := r.Int63() % int64(len(letters))
		value[i] = letters[int(r)]
	}
	return &ConfirmCode{value: value}
}

// 文字列の確認コードを返す
// @return
// 確認コード
func (cc *ConfirmCode) String() string {
	return string(cc.value)
}
