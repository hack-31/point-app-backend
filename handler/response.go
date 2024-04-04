package handler

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/hack-31/point-app-backend/middleware"
)

// クライアントへの返すレスポンスの型
type Responses struct {
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

// エラーのレスポンスの型
type ErrRes struct {
	StatusCode int    `json:"statusCode"`
	Title      string `json:"title"`
	Message    string `json:"message"`
}

// APIレスポンスの作成（成功時）
//
// @params
// ctx ginのコンテキスト
// StatusCode ステータスコード
// Message メッセージ
// Data 返却するデータ
func APIResponse(ctx *gin.Context, StatusCode int, Message string, Data interface{}) {
	jsonResponse := &Responses{
		StatusCode: StatusCode,
		Message:    Message,
		Data:       Data,
	}
	ctx.JSON(StatusCode, jsonResponse)

	// ログに使う情報登録
	js, err := json.Marshal(jsonResponse)
	if err != nil {
		panic("can not marshal")
	}
	middleware.NewLog().
		Rsp(js).
		Logging(ctx)

}

// エラーレスポンス作成
//
// @params
// ctx ginのコンテキスト
// StatusCode ステータスコード
// Title エラータイトル
// Message メッセージ
func ErrResponse(ctx *gin.Context, StatusCode int, Title, Message string, Err error) {
	res := ErrRes{
		StatusCode: StatusCode,
		Title:      Title,
		Message:    Message,
	}
	defer ctx.AbortWithStatus(StatusCode)
	ctx.JSON(StatusCode, res)

	// ログに使う情報登録
	js, err := json.Marshal(res)
	if err != nil {
		panic("can not marshal")
	}
	middleware.NewLog().
		ErrCode(StatusCode).
		Err(Err).
		Rsp(js).
		Logging(ctx)
}
