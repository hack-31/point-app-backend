package repository

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hack-31/point-app-backend/domain/model"
	"github.com/hack-31/point-app-backend/repository/testdata"
	"github.com/hack-31/point-app-backend/testutil"
	"github.com/stretchr/testify/assert"
)

func TestRepository_GetAquistionPoint(t *testing.T) {
	t.Parallel()
	type want struct {
		users map[model.UserID]int
	}
	type input struct {
		userIDs []model.UserID
	}

	tests := map[string]struct {
		input input
		want  want
	}{
		"指定ユーザーの獲得ポイントを取得する": {
			input: input{
				userIDs: []model.UserID{1, 2, 3},
			},
			want: want{
				users: map[model.UserID]int{
					1: 1200,
					2: 500,
					3: 300,
				},
			},
		},
	}

	for n, tt := range tests {
		tt := tt
		t.Run(n, func(t *testing.T) {
			t.Parallel()
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			tx, err := testutil.OpenDBForTest(t).BeginTxx(ctx, nil)
			assert.NoError(t, err)

			t.Cleanup(func() {
				err := tx.Rollback()
				assert.NoError(t, err)
			})

			testdata.Transactions(t, ctx, tx, func(users model.Transactions) {})

			// 実行
			r := &Repository{}
			gots, err := r.GetAquistionPoint(ctx, tx, tt.input.userIDs)

			// アサーション
			assert.Nil(t, err)
			assert.Equal(t, tt.want.users, gots)
		})
	}
}
