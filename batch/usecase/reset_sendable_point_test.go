package usecase

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	mock_domain "github.com/hack-31/point-app-backend/domain/_mock"
	"github.com/hack-31/point-app-backend/repository"
	mock_repository "github.com/hack-31/point-app-backend/repository/_mock"
	"github.com/hack-31/point-app-backend/testutil"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestResetPoint(t *testing.T) {
	t.Parallel()
	type input struct {
		point int
	}
	type updateAllSendablePointInput struct {
		point int
	}
	type updateAllSendablePointOutput struct {
		err error
	}
	type want struct {
		err error
	}
	tests := map[string]struct {
		input                        input
		updateAllSendablePointInput  updateAllSendablePointInput
		updateAllSendablePointOutput updateAllSendablePointOutput
		want                         want
	}{
		"送付可能ポイント1000ポイントで更新できる": {
			input: input{
				point: 1000,
			},
			updateAllSendablePointInput: updateAllSendablePointInput{
				point: 1000,
			},
			updateAllSendablePointOutput: updateAllSendablePointOutput{
				err: nil,
			},
			want: want{
				err: nil,
			},
		},
		"更新でエラー": {
			input: input{
				point: 1000,
			},
			updateAllSendablePointInput: updateAllSendablePointInput{
				point: 1000,
			},
			updateAllSendablePointOutput: updateAllSendablePointOutput{
				err: repository.ErrDBException,
			},
			want: want{
				err: repository.ErrDBException,
			},
		},
	}

	for n, tt := range tests {
		tt := tt
		t.Run(n, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			t.Cleanup(func() { ctrl.Finish() })

			// コンテキストの設定
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			// トランザクションのモック設定
			mockTx := testutil.NewTxForMock(t, ctx)

			mockBeginner := mock_repository.NewMockBeginner(ctrl)
			mockBeginner.
				EXPECT().
				BeginTxx(ctx, nil).
				Return(mockTx, nil)

			mockPintRepo := mock_domain.NewMockPointRepo(ctrl)
			mockPintRepo.
				EXPECT().
				UpdateAllSendablePoint(ctx, mockTx, tt.updateAllSendablePointInput.point).
				Return(tt.updateAllSendablePointOutput.err)

			r := ResetSendablePoint{
				Tx:        mockBeginner,
				PointRepo: mockPintRepo,
			}
			gots := r.ResetPoint(ctx, tt.input.point)
			assert.Equal(t, tt.want.err, gots)
		})
	}
}
