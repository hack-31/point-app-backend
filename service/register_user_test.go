package service

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hack-31/point-app-backend/repository"
	"github.com/hack-31/point-app-backend/repository/entity"
	"github.com/stretchr/testify/assert"
)

func TestRegisterUser(t *testing.T) {
	t.Parallel()

	wantConfirmCode := "1234"
	wantTemporaryUserID := "temp-user-id"
	wantUserID := entity.UserID(10)
	wantUser := &entity.User{
		ID:             wantUserID,
		FirstName:      "太郎",
		FirstNameKana:  "たろう",
		FamilyName:     "山田",
		FamilyNameKana: "やまだ",
		Password:       "pass",
		Email:          "test@sample.com",
		SendingPoint:   1000,
	}
	wantToken := "token"
	wantTemporaryUser := "太郎\nたろう\n山田\nやまだ\ntest@sample.com\npass"

	// 入力
	type input struct {
		confirmCode     string
		temporaryUserId string
	}
	// 結果
	type want struct {
		user  *entity.User
		token string
		err   error
	}
	type UserAdderMockParameter struct {
		in  *entity.User
		err error
	}
	type tempUserMockParameter struct {
		rps string
		err error
	}
	type generateTokenMockParameter struct {
		rsp []byte
		err error
	}

	tests := map[string]struct {
		wants want
		input input
		// ユーザをDBに保存するリポジトリのモックパラメータ
		uaprm UserAdderMockParameter
		// 一時ユーザを取得するCacheのモックパラメータ
		tumprm tempUserMockParameter
		// JWTを発行するモックパラメータ
		gtmprm generateTokenMockParameter
	}{
		"確認コードとtemmporaryUserIDが正しい場合は登録ユーザを返す": {
			input: input{
				confirmCode:     wantConfirmCode,
				temporaryUserId: wantTemporaryUserID,
			},
			wants: want{
				err:   nil,
				user:  wantUser,
				token: wantToken,
			},
			uaprm: UserAdderMockParameter{
				in:  wantUser,
				err: nil,
			},
			tumprm: tempUserMockParameter{
				rps: wantTemporaryUser,
				err: nil,
			},
			gtmprm: generateTokenMockParameter{
				rsp: []byte(wantToken),
				err: nil,
			},
		},
		"Cacheより一時保存したユーザ情報が取得できない場合は、エラー": {
			input: input{
				confirmCode:     wantConfirmCode,
				temporaryUserId: wantTemporaryUserID,
			},
			wants: want{
				err:   repository.ErrNotFoundSession,
				user:  nil,
				token: "",
			},
			uaprm: UserAdderMockParameter{
				in:  nil,
				err: nil,
			},
			tumprm: tempUserMockParameter{
				rps: "",
				err: repository.ErrNotFoundSession,
			},
			gtmprm: generateTokenMockParameter{
				rsp: []byte(wantToken),
				err: nil,
			},
		},
		"ユーザ情報をDBに保存する際、メールアドレスが重複し登録できない時は、エラー": {
			input: input{
				confirmCode:     wantConfirmCode,
				temporaryUserId: wantTemporaryUserID,
			},
			wants: want{
				err:   repository.ErrAlreadyEntry,
				user:  nil,
				token: "",
			},
			uaprm: UserAdderMockParameter{
				in:  wantUser,
				err: repository.ErrAlreadyEntry,
			},
			tumprm: tempUserMockParameter{
				rps: wantTemporaryUser,
				err: nil,
			},
			gtmprm: generateTokenMockParameter{
				rsp: []byte(wantToken),
				err: nil,
			},
		},
	}

	for n, tt := range tests {
		tt := tt
		t.Run(n, func(t *testing.T) {
			t.Parallel()
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Set("email", "test@sample.com")

			// モックの定義
			moqDB := &ExecerMock{}
			moqRepo := &UserRepoMock{}
			moqCache := &CacheMock{}
			moqCache.LoadFunc = func(pctx context.Context, key string) (string, error) {
				assert.Equal(t, ctx, pctx)
				return tt.tumprm.rps, tt.tumprm.err
			}
			moqCache.DeleteFunc = func(pctx context.Context, key string) error {
				return nil
			}
			moqTokenGenerator := &TokenGeneratorMock{}
			moqTokenGenerator.GenerateTokenFunc = func(pctx context.Context, u entity.User) ([]byte, error) {
				assert.Equal(t, *wantUser, u)
				return tt.gtmprm.rsp, tt.gtmprm.err
			}
			moqRepo.RegisterUserFunc = func(pctx context.Context, db repository.Execer, user *entity.User) error {
				user.ID = wantUserID
				assert.Equal(t, tt.uaprm.in, user)
				assert.Equal(t, ctx, pctx)
				assert.Equal(t, moqDB, db)
				return tt.uaprm.err
			}
			ru := &RegisterUser{
				DB:             moqDB,
				Cache:          moqCache,
				Repo:           moqRepo,
				TokenGenerator: moqTokenGenerator,
			}

			// サービスの実行
			gotUser, gotToken, gotErr := ru.RegisterUser(ctx, tt.input.temporaryUserId, tt.input.confirmCode)

			// アサーション
			assert.ErrorIs(t, gotErr, tt.wants.err)
			assert.Equal(t, tt.wants.user, gotUser)
			assert.Equal(t, tt.wants.token, gotToken)
		})
	}
}
