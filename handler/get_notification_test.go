package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hack-31/point-app-backend/repository"
	"github.com/hack-31/point-app-backend/repository/entity"
	"github.com/hack-31/point-app-backend/service"
	"github.com/hack-31/point-app-backend/testutil"
)

func TestGetNotification(t *testing.T) {
	t.Parallel()
	type want struct {
		status  int
		rspFile string
	}

	type getNotification struct {
		err error
		res service.GetNotificationResponse
	}
	type input struct {
		id string
	}

	tests := map[string]struct {
		input           input
		want            want
		getNotification getNotification
	}{
		"パスパラメータのidに対するお知らせ詳細と200を返す": {
			want: want{
				status:  http.StatusOK,
				rspFile: "testdata/get_notification/200_rsp.json.golden",
			},
			input: input{
				id: "1",
			},
			getNotification: getNotification{
				err: nil,
				res: service.GetNotificationResponse{
					ID:          1,
					Title:       "ポイント送付のお知らせ",
					Description: "斉藤さんよりポイントが100ポイント送付されました。",
					IsChecked:   true,
					CreatedAt:   "2022/12/08 11:08:08",
				},
			},
		},
		"お知らせIDに対するお知らせ詳細データが存在しない場合は、404エラーを返す": {
			input: input{
				id: "1",
			},
			want: want{
				status:  http.StatusNotFound,
				rspFile: "testdata/get_notification/404_rsp.json.golden",
			},
			getNotification: getNotification{
				err: repository.ErrNotFound,
				res: service.GetNotificationResponse{},
			},
		},
		"パスパラメータidに文字列「test」の場合は400エラーを返す": {
			input: input{
				id: "test",
			},
			want: want{
				status:  http.StatusBadRequest,
				rspFile: "testdata/get_notification/400_rsp.json.golden",
			},
			getNotification: getNotification{
				err: nil,
				res: service.GetNotificationResponse{},
			},
		},
	}

	for n, tt := range tests {
		tt := tt
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			// サービス層のモック定義
			moq := &GetNotificationServiceMock{
				GetNotificationFunc: func(ctx *gin.Context, notificationID entity.NotificationID) (service.GetNotificationResponse, error) {
					return tt.getNotification.res, tt.getNotification.err
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
			gn := NewGetNotification(moq)
			gn.ServeHTTP(c)

			// レスポンス
			resp := w.Result()
			// 検証
			testutil.AssertResponse(t, resp, tt.want.status, testutil.LoadFile(t, tt.want.rspFile))
		})
	}
}
