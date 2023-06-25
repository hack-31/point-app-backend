package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hack-31/point-app-backend/domain/model"
	"github.com/hack-31/point-app-backend/service"
	"github.com/hack-31/point-app-backend/testutil"
)

func TestGetNotifications(t *testing.T) {
	type want struct {
		status  int
		rspFile string
	}
	type getNotifications struct {
		err error
		res service.GetNotificationsResponse
	}
	type input struct {
		size      string
		nextToken string
	}
	tests := map[string]struct {
		input            input
		want             want
		getNotifications getNotifications
	}{
		"クエリのsize、nextTokenに応じてお知らせ一覧と200を返す": {
			want: want{
				status:  http.StatusOK,
				rspFile: "testdata/get_notifications/200_rsp.json.golden",
			},
			input: input{
				size:      "5",
				nextToken: "10",
			},
			getNotifications: getNotifications{
				err: nil,
				res: service.GetNotificationsResponse{
					NextToken: "100",
					Notifications: []struct {
						ID          model.NotificationID
						Title       string
						Description string
						IsChecked   bool
						CreatedAt   string
					}{
						{ID: 1, Title: "ポイント送付のお知らせ", Description: "太郎から100ポイント送付されました。", IsChecked: false},
					},
				},
			},
		},
		"sizeが文字列の場合は400エラーを返す": {
			want: want{
				status:  http.StatusBadRequest,
				rspFile: "testdata/get_notifications/400_rsp.json.golden",
			},
			input: input{
				size:      "test",
				nextToken: "10",
			},
			getNotifications: getNotifications{
				err: nil,
				res: service.GetNotificationsResponse{
					NextToken: "0",
					Notifications: []struct {
						ID          model.NotificationID
						Title       string
						Description string
						IsChecked   bool
						CreatedAt   string
					}{},
				},
			},
		},
	}

	for n, tt := range tests {
		tt := tt
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			// サービス層のモック定義
			moq := &GetNotificationsServiceMock{
				GetNotificationsFunc: func(ctx *gin.Context, nextToken, size string) (service.GetNotificationsResponse, error) {
					return tt.getNotifications.res, tt.getNotifications.err
				},
			}

			// テストデータを挿入
			// コンテキストの作成
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// クエリパラメータの設定
			url := fmt.Sprintf("/notifications?size=%s&nextToken=%s", tt.input.size, tt.input.nextToken)
			c.Request, _ = http.NewRequest("GET", url, nil)

			// リクエスト送信
			gn := NewGetNotifications(moq)
			gn.ServeHTTP(c)

			// レスポンス
			resp := w.Result()
			// 検証
			testutil.AssertResponse(t, resp, tt.want.status, testutil.LoadFile(t, tt.want.rspFile))
		})
	}
}
