package auth

import (
	"context"
	_ "embed"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hack-31/point-app-backend/constant"
	"github.com/hack-31/point-app-backend/entity"
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
		return nil, fmt.Errorf("failed in NewJWTer: private key: %w", err)
	}
	pubkey, err := parse(rawPubKey)
	if err != nil {
		return nil, fmt.Errorf("failed in NewJWTer: public key: %w", err)
	}
	j.PrivateKey = privkey
	j.PublicKey = pubkey
	j.Clocker = c
	return j, nil
}

func parse(rawKey []byte) (jwk.Key, error) {
	key, err := jwk.ParseKey(rawKey, jwk.WithPEM(true))
	if err != nil {
		return nil, err
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
func (j *JWTer) GenerateToken(ctx context.Context, u entity.User) ([]byte, error) {
	tok, err := jwt.NewBuilder().
		JwtID(uuid.New().String()).
		Issuer(`github.com/hack-31/point-app-backend`).
		Subject("access_token").
		IssuedAt(j.Clocker.Now()).
		Expiration(j.Clocker.Now().Add(time.Duration(constant.MaxTokenExpiration_m)*time.Minute)).
		Claim(Email, u.Email).
		Claim(UserID, u.ID).
		Build()
	if err != nil {
		return nil, fmt.Errorf("GenerateToken: failed to build token: %w", err)
	}
	if err := j.Store.Save(ctx, fmt.Sprint(u.ID), tok.JwtID(), time.Duration(constant.TokenExpiration_m)); err != nil {
		return nil, err
	}

	// 署名
	signed, err := jwt.Sign(tok, jwt.WithKey(jwa.RS256, j.PrivateKey))
	if err != nil {
		return nil, err
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
		return nil, err
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
		return err
	}

	// 有効期限が切れていないか確認
	if err := jwt.Validate(token, jwt.WithClock(j.Clocker)); err != nil {
		return fmt.Errorf("GetToken: failed to validate token: %w", err)
	}

	// キャッシュに対しても有効期限を確認
	id, ok := token.Get(UserID)
	if !ok {
		return fmt.Errorf("not found %s", UserID)
	}
	uid := fmt.Sprintf("%v", id)
	jwi, err := j.Store.Load(ctx, uid)
	if  err != nil {
		return fmt.Errorf("GetToken: %v expired: %w", id, err)
	}

	// 他のログインを検査
	if jwi != token.JwtID() {
		return fmt.Errorf("expired token %s because login another", jwi)
	}

	// 有効期限延長
	j.Store.Expire(ctx, uid, time.Duration(constant.TokenExpiration_m))

	// コンテキストにユーザ情報追加
	ctx.Set(UserID, uid)
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
