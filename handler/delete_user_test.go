package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hack-31/point-app-backend/domain/model"
	"github.com/hack-31/point-app-backend/repository"
	"github.com/hack-31/point-app-backend/testutil"
	"github.com/stretchr/testify/assert"
)

func TestDeleteUser(t *testing.T) {
	t.Parallel()

	type want struct {
		status  int
		rspFile string
	}

	type deleteUserInput struct {
		userID model.UserID
	}
	type deleteUser struct {
		input  deleteUserInput
		output error
	}

	type input struct {
		id string
	}

	tests := map[string]struct {
		input      input
		deleteUser deleteUser
		want       want
	}{
		"パスパラメータのidに対するお知らせ詳細と201を返す": {
			input: input{
				id: "1",
			},
			deleteUser: deleteUser{
				input: deleteUserInput{
					userID: 1,
				},
				output: nil,
			},
			want: want{
				status:  http.StatusCreated,
				rspFile: "testdata/delete_user/201_rsp.json.golden",
			},
		},
		"id=-1を指定した場合は400を返す": {
			input: input{
				id: "-1",
			},
			deleteUser: deleteUser{
				input: deleteUserInput{
					userID: 1,
				},
				output: nil,
			},
			want: want{
				status:  http.StatusBadRequest,
				rspFile: "testdata/delete_user/400_rsp.json.golden",
			},
		},
		"ユーザーが存在しないエラーの時404を返す": {
			input: input{
				id: "1",
			},
			deleteUser: deleteUser{
				input: deleteUserInput{
					userID: 1,
				},
				output: repository.ErrNotUser,
			},
			want: want{
				status:  http.StatusNotFound,
				rspFile: "testdata/delete_user/404_rsp.json.golden",
			},
		},
	}

	for n, tt := range tests {
		tt := tt
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			// サービス層のモック定義
			moq := &DeleteUserServiceMock{
				DeleteUserFunc: func(ctx *gin.Context, userID model.UserID) error {
					assert.Equal(t, tt.deleteUser.input.userID, userID)
					return tt.deleteUser.output
				},
			}

			// テストデータを挿入
			// コンテキストの作成
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// パスパラメータの設定
			param := gin.Param{Key: "id", Value: tt.input.id}
			c.Params = gin.Params{param}

			// リクエスト送信
			du := NewDeleteUser(moq)
			du.ServeHTTP(c)

			// レスポンス
			resp := w.Result()
			// 検証
			testutil.AssertResponse(t, resp, tt.want.status, testutil.LoadFile(t, tt.want.rspFile))
		})
	}
}
