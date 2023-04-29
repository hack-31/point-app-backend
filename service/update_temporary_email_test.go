package service

import (
	"context"
	"database/sql"
	"errors"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"
	"github.com/hack-31/point-app-backend/domain/model"
	"github.com/hack-31/point-app-backend/repository"
)

func TestUpdateTemporaryEmail(t *testing.T) {
	inputEmail := "yamada.tarotestetst@sample.com"
	inputExistEmail := "yamada@sample.com"
	wantUserID := model.UserID(10)
	wantUser := &model.User{
		ID:             wantUserID,
		FirstName:      "太郎",
		FirstNameKana:  "たろう",
		FamilyName:     "山田",
		FamilyNameKana: "やまだ",
		Password:       "pass",
		Email:          "yamada@sample.com",
		SendingPoint:   1000,
	}

	// 入力
	type input struct {
		email string
	}
	// 結果
	type want struct {
		temporaryEmailId string
		err              error
	}

	tests := map[string]struct {
		wants want
		input input
	}{
		"すでに登録されているメールアドレスである場合は、エラーを返す": {
			input: input{
				email: inputExistEmail,
			},
			wants: want{
				temporaryEmailId: "",
				err:              repository.ErrAlreadyEntry,
			},
		},
	}

	for n, tt := range tests {
		tt := tt
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			// モックの定義
			moqDB := &QueryerMock{}
			moqCache := &CacheMock{}
			moqRepo := &UserRepoMock{}
			moqCache.SaveFunc = func(ctx context.Context, key, value string, minute time.Duration) error {
				return nil
			}
			moqRepo.FindUserByEmailFunc = func(ctx context.Context, db repository.Queryer, email *string) (model.User, error) {
				if tt.input.email == inputEmail {
					return model.User{}, sql.ErrNoRows
				}
				if tt.input.email == inputExistEmail {
					return *wantUser, nil
				}
				return *wantUser, nil
			}
			ute := &UpdateTemporaryEmail{
				DB:    moqDB,
				Cache: moqCache,
				Repo:  moqRepo,
			}

			// サービス実行
			gotTemporaryEmailId, gotErr := ute.UpdateTemporaryEmail(ctx, tt.input.email)

			// アサーション
			if !errors.Is(gotErr, tt.wants.err) {
				t.Errorf("%s error\nwant: %+v\ngot: %+v\n", t.Name(), gotErr, tt.wants.err)
			}
			if d := cmp.Diff(gotTemporaryEmailId, tt.wants.temporaryEmailId); len(d) != 0 {
				t.Errorf("differs: (-got +want)\n%s", d)
			}
		})
	}
}
