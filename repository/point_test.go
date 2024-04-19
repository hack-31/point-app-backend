package repository

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hack-31/point-app-backend/repository/entity"
	"github.com/hack-31/point-app-backend/repository/testdata"
	"github.com/hack-31/point-app-backend/testutil"
	"github.com/stretchr/testify/assert"
)

func TestRepository_users_UpdateSendablePoint(t *testing.T) {
	type want struct {
		afterPoint int
		err        error
	}
	type input struct {
		beforePoint int
		id          entity.UserID
	}

	tests := map[string]struct {
		input input
		want  want
	}{
		"ユーザーID_1のメールを更新する": {
			input: input{
				beforePoint: 1000,
				id:          1,
			},
			want: want{
				afterPoint: 100,
				err:        nil,
			},
		},
	}

	for n, tt := range tests {
		t.Run(n, func(t *testing.T) {
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			tx, err := testutil.OpenDBForTest(t).BeginTxx(ctx, nil)
			assert.NoError(t, err)
			t.Cleanup(func() {
				err := tx.Rollback()
				assert.NoError(t, err)
			})

			testdata.Users(t, ctx, tx, func(users entity.Users) {
				users[tt.input.id-1].SendingPoint = tt.input.beforePoint
			})

			// 実行
			r := &Repository{}
			err = r.UpdateSendablePoint(ctx, tx, tt.input.id, tt.want.afterPoint)
			assert.Nil(t, err)

			// 確認
			gots, err := r.GetUserByID(ctx, tx, tt.input.id)
			assert.Nil(t, err)
			assert.Equal(t, tt.want.afterPoint, gots.SendingPoint)
		})
	}
}
