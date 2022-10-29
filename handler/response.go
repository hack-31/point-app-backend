package handler

import "github.com/gin-gonic/gin"

// クライアントへの返すレスポンスの型
type Responses struct {
	StatusCode int         `json:"statusCode"`
	Method     string      `json:"method"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

// APIレスポンスの作成
// @params ctx ginのコンテキスト
//
// @params Message メッセージ
//
// @params StatusCode ステータスコード
//
// @params Method HTTPメソッド
//
// @params Data 返却するデータ
func APIResponse(ctx *gin.Context, Message string, StatusCode int, Method string, Data interface{}) {
	jsonResponse := Responses{
		StatusCode: StatusCode,
		Method:     Method,
		Message:    Message,
		Data:       Data,
	}

	if StatusCode >= 400 {
		ctx.JSON(StatusCode, jsonResponse)
		defer ctx.AbortWithStatus(StatusCode)
		return
	}
	ctx.JSON(StatusCode, jsonResponse)
}
