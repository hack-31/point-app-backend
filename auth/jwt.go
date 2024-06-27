package auth

import (
	"context"
	_ "embed"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hack-31/point-app-backend/constant"
	"github.com/hack-31/point-app-backend/domain/model"
	"github.com/hack-31/point-app-backend/repository/entities"
	"github.com/hack-31/point-app-backend/utils/clock"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

const (
	UserID = "user_id"
	Email  = "email"
)

//go:embed certificate/secret.pem
var rawPrivKey []byte

//go:embed certificate/public.pem
var rawPubKey []byte

type JWTer struct {
	PrivateKey, PublicKey jwk.Key
	Store                 Store
	Clocker               clock.Clocker
}

// jWTのインスタンス
func NewJWTer(s Store, c clock.Clocker) (*JWTer, error) {
	j := &JWTer{Store: s}
	privkey, err := parse(rawPrivKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse private key")
	}
	pubkey, err := parse(rawPubKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse public key")
	}
	j.PrivateKey = privkey
	j.PublicKey = pubkey
	j.Clocker = c
	return j, nil
}

func parse(rawKey []byte) (jwk.Key, error) {
	key, err := jwk.ParseKey(rawKey, jwk.WithPEM(true))
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse key")
	}
	return key, nil
}

// アクセストークンの作成
// @params
// ctx コンテキスト
// u ユーザエンティティ
//
// @returns
// token アクセストークン
func (j *JWTer) GenerateToken(ctx context.Context, u entities.User) ([]byte, error) {
	tok, err := jwt.NewBuilder().
		JwtID(uuid.New().String()).
		IssuedAt(j.Clocker.Now()).
		Expiration(j.Clocker.Now().Add(time.Duration(constant.MaxTokenExpiration_m)*time.Minute)).
		Claim(Email, u.Email).
		Claim(UserID, u.ID).
		Build()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build token")
	}
	if err := j.Store.Save(ctx, fmt.Sprint(u.ID), tok.JwtID(), time.Duration(constant.TokenExpiration_m)); err != nil {
		return nil, errors.Wrap(err, "failed to save token")
	}

	// 署名
	signed, err := jwt.Sign(tok, jwt.WithKey(jwa.RS256, j.PrivateKey))
	if err != nil {
		return nil, errors.Wrap(err, "failed to sign token")
	}
	return signed, nil
}

// トークンを取得し解析
// @params
// ctx コンテキスト
// r リクエスト情報
//
// @returns
// トークン
func (j *JWTer) GetToken(ctx context.Context, r *http.Request) (jwt.Token, error) {
	token, err := jwt.ParseRequest(
		r,
		jwt.WithKey(jwa.RS256, j.PublicKey),
		jwt.WithValidate(false),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse token")
	}
	return token, nil
}

// トークンを解析し、contextにuserIDとEmailをセットする
// @params
// ctx コンテキスト
func (j *JWTer) FillContext(ctx *gin.Context) error {
	// トークンを解析
	token, err := j.GetToken(ctx.Request.Context(), ctx.Request)
	if err != nil {
		return errors.Wrap(err, "failed to get token")
	}

	// 有効期限が切れていないか確認
	if err := jwt.Validate(token, jwt.WithClock(j.Clocker)); err != nil {
		return errors.Wrap(err, "failed to validate token")
	}

	// キャッシュに対しても有効期限を確認
	id, ok := token.Get(UserID)
	if !ok {
		return errors.Wrapf(err, "failed to get %s", UserID)
	}
	uid := fmt.Sprintf("%v", id)
	jwi, err := j.Store.Load(ctx, uid)
	if err != nil {
		return errors.Wrapf(err, "failed to load %s", id)
	}

	// 他のログインを検査
	if jwi != token.JwtID() {
		return errors.Wrapf(err, "expired token %s because login another", jwi)
	}

	// 有効期限延長
	if err = j.Store.Expire(ctx, uid, time.Duration(constant.TokenExpiration_m)); err != nil {
		return errors.Wrap(err, "can not be extended")
	}

	// コンテキストにユーザ情報追加
	intUid, _ := strconv.ParseInt(uid, 10, 64)
	ctx.Set(UserID, model.UserID(intUid))
	SetEmail(ctx, token)

	return nil
}

// メールをコンテキストに代入
//
// @paramss
// ctx コンテキスト
// tok トークン
func SetEmail(ctx *gin.Context, tok jwt.Token) {
	get, ok := tok.Get(Email)
	if !ok {
		ctx.Set(Email, "")
		return
	}
	ctx.Set(Email, get)
}
