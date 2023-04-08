package handler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hack-31/point-app-backend/domain/model"
	"github.com/hack-31/point-app-backend/repository"
	"github.com/hack-31/point-app-backend/testutil"
	"github.com/hack-31/point-app-backend/utils/clock"
)

func TestRegisterUser(t *testing.T) {
	type want struct {
		status  int
		rspFile string
	}

	tests := map[string]struct {
		reqFile string
		want    want
	}{
		"正しいリクエストの時は201となる": {
			reqFile: "testdata/register_user/201_req.json.golden",
			want: want{
				status:  http.StatusCreated,
				rspFile: "testdata/register_user/201_rsp.json.golden",
			},
		},
		"登録済みのユーザは409エラーを返す": {
			reqFile: "testdata/register_user/409_req.json.golden",
			want: want{
				status:  http.StatusConflict,
				rspFile: "testdata/register_user/409_rsp.json.golden",
			},
		},
		"確認コードが無効の場合は401エラーを返す": {
			reqFile: "testdata/register_user/401_req.json.golden",
			want: want{
				status:  http.StatusUnauthorized,
				rspFile: "testdata/register_user/401_rsp.json.golden",
			},
		},
	}

	for n, tt := range tests {
		tt := tt
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			// サービス層のモック定義
			moq := &RegisterUserServiceMock{}
			moq.RegisterUserFunc = func(ctx context.Context, temporaryUserId, confirmCode string) (*model.User, string, error) {
				clocker := clock.FixedClocker{}
				if tt.want.status == http.StatusCreated {
					return &model.User{
							ID:             24,
							FirstName:      "山田",
							FirstNameKana:  "やまだ",
							FamilyName:     "太郎",
							FamilyNameKana: "たろう",
							Email:          "hoge@hoge.com",
							CreatedAt:      clocker.Now(),
							UpdateAt:       clocker.Now()},
						"json-web-token",
						nil
				}
				if tt.want.status == http.StatusConflict {
					return nil, "", repository.ErrAlreadyEntry
				}
				if tt.want.status == http.StatusUnauthorized {
					return nil, "", repository.ErrNotFoundSession
				}

				return nil, "", errors.New("error from mock")
			}

			// テストデータを挿入
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest("POST", "/users", bytes.NewReader(testutil.LoadFile(t, tt.reqFile)))
			ru := RegisterUser{
				Service: moq,
			}

			// リクエスト送信
			ru.ServeHTTP(c)

			// レスポンス
			resp := w.Result()
			// 検証
			testutil.AssertResponse(t, resp, tt.want.status, testutil.LoadFile(t, tt.want.rspFile))
		})
	}
}
