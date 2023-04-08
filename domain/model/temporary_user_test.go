package model_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hack-31/point-app-backend/domain/model"
)

func TestTemporaryUser_Join(t *testing.T) {
	t.Parallel()

	// 入力
	type input struct {
		firstName, firstNameKana, familyName, familyNameKana, email, password string
	}

	tests := map[string]struct {
		want  string
		input input
	}{
		"ユーザ情報をRedisに一時保存するため、改行文字で区切って一つの文字列にする": {
			input: input{
				firstName:      "太郎",
				firstNameKana:  "たろう",
				familyName:     "山田",
				familyNameKana: "やまだ",
				email:          "yamada.taro@example.com",
				password:       "pass",
			},
			want: "太郎\nたろう\n山田\nやまだ\nyamada.taro@example.com\npass",
		},
	}

	for n, tt := range tests {
		tt := tt
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			// 実行
			tempUser := model.NewTemporaryUserString("")
			got := tempUser.Join(
				tt.input.firstName,
				tt.input.firstNameKana,
				tt.input.familyName,
				tt.input.familyNameKana,
				tt.input.email,
				tt.input.password,
			)

			// アサーション
			if d := cmp.Diff(got, tt.want); len(d) != 0 {
				t.Errorf("%s error\nwant:\n%+v\n\ngot:\n%+v\n", t.Name(), got, tt.want)
			}
		})
	}
}

func TestTemporaryUser_Split(t *testing.T) {
	t.Parallel()

	type want struct {
		firstName, firstNameKana, familyName, familyNameKana, email, password string
	}

	tests := map[string]struct {
		input string
		want  want
	}{
		"一時保存されたユーザ情報文字列を分割する": {
			input: "太郎\nたろう\n山田\nやまだ\nyamada.taro@example.com\npass",
			want: want{
				firstName:      "太郎",
				firstNameKana:  "たろう",
				familyName:     "山田",
				familyNameKana: "やまだ",
				email:          "yamada.taro@example.com",
				password:       "pass",
			},
		},
	}

	for n, tt := range tests {
		tt := tt
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			// 実行
			tempUser := model.NewTemporaryUserString(tt.input)
			gotFirstName, gotFirstNameKana, gotFamilyName, gotFamilyNameKana, gotEmail, gotPassword := tempUser.Split()

			// アサーション
			if d := cmp.Diff(gotFirstName, tt.want.firstName); len(d) != 0 {
				t.Errorf("differs: (-got +want)\n%s", d)
			}
			if d := cmp.Diff(gotFirstNameKana, tt.want.firstNameKana); len(d) != 0 {
				t.Errorf("differs: (-got +want)\n%s", d)
			}
			if d := cmp.Diff(gotFamilyName, tt.want.familyName); len(d) != 0 {
				t.Errorf("differs: (-got +want)\n%s", d)
			}
			if d := cmp.Diff(gotFamilyNameKana, tt.want.familyNameKana); len(d) != 0 {
				t.Errorf("differs: (-got +want)\n%s", d)
			}
			if d := cmp.Diff(gotPassword, tt.want.password); len(d) != 0 {
				t.Errorf("differs: (-got +want)\n%s", d)
			}
			if d := cmp.Diff(gotEmail, tt.want.email); len(d) != 0 {
				t.Errorf("differs: (-got +want)\n%s", d)
			}
		})
	}
}
