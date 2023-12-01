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

	notificationCntChan, err := gunc.Service.GetUncheckedNotificationCounts(ctx)
	if err != nil {
		ctx.SSEvent("error", err.Error())
		ctx.Writer.Flush()
		return
	}

	// お知らせ通知を待機
	for {
		select {
		case <-ctx.Request.Context().Done():
			ctx.SSEvent("error", "キャンセルされました。")
			ctx.Writer.Flush()
			return
		case <-ctx.Done():
			ctx.SSEvent("error", "キャンセルされました。")
			ctx.Writer.Flush()
			return
		case cnt, ok := <-notificationCntChan:
			if !ok {
				ctx.SSEvent("error", "キャンセルされました。")
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
		}
	}
}
