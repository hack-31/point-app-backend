package handler

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
)

type GetUncheckedNotificationCount struct {
	Service GetUncheckedNotificationCountService
}

func NewGetUncheckedNotificationCount(s GetUncheckedNotificationCountService) *GetUncheckedNotificationCount {
	return &GetUncheckedNotificationCount{Service: s}
}

// お知らせ数取得ハンドラー
//
// @param ctx ginContext
func (gunc *GetUncheckedNotificationCount) ServeHTTP(ctx *gin.Context) {
	// ヘッダー情報設定
	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")

	notificationCntChan := make(chan int)
	defer close(notificationCntChan)

	//　初回お知らせ数の返却
	cnt, err := gunc.Service.GetUncheckedNotificationCount(ctx, notificationCntChan)
	if err != nil {
		ctx.SSEvent("error", err.Error())
		ctx.Writer.Flush()
		return
	}
	res := struct {
		Count int `json:"count"`
	}{
		Count: cnt,
	}
	jsonData, err := json.Marshal(res)
	if err != nil {
		ctx.SSEvent("error", err.Error())
		ctx.Writer.Flush()
		return
	}
	ctx.SSEvent("message", string(jsonData))
	ctx.Writer.Flush()

	// お知らせ通知を待機
	for {
		select {
		case <-ctx.Request.Context().Done():
			return
		case cnt := <-notificationCntChan:
			res := struct {
				Count int `json:"count"`
			}{
				Count: cnt,
			}
			jsonData, err := json.Marshal(res)
			if err != nil {
				ctx.SSEvent("error", err.Error())
				ctx.Writer.Flush()
				return
			}
			ctx.SSEvent("message", string(jsonData))
			ctx.Writer.Flush()
		}
	}
}
