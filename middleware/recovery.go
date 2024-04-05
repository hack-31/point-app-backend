package middleware

import (
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hack-31/point-app-backend/utils/clock"
	"github.com/rs/zerolog"
)

// Recovery は、リクエストのパニックを回復する
func Recovery(ctx *gin.Context) {
	ctx.Set("startTime", time.Now())
	defer func() {
		if err := recover(); err != nil {
			// パスの取得
			path := ctx.Request.URL.Path
			raw := ctx.Request.URL.RawQuery
			if raw != "" {
				path = path + "?" + raw
			}

			// エラーコード取得
			gec, exist := ctx.Get(errCodeKey)
			var ec int
			if exist {
				ec = gec.(int)
			}

			// 経過時刻
			gst, exists := ctx.Get("startTime")
			var startTime time.Time
			if exists {
				startTime = gst.(time.Time)
			}
			elapsedTime := time.Since(startTime)

			// ログ出力
			logger := zerolog.New(os.Stdout).
				With().
				Timestamp().
				Logger().
				Output(zerolog.ConsoleWriter{Out: os.Stderr, FormatTimestamp: func(i any) string {
					clocker := clock.RealClocker{}
					// 日本時間にする
					jst := time.FixedZone("Asia/Tokyo", 9*60*60)
					return clocker.Now().In(jst).Format(time.RFC3339)
				}})

			logger.Error().
				// レスポンス
				Int("status", http.StatusInternalServerError).
				// リクエスト
				Str("URI", path).
				Str("method", ctx.Request.Method).
				Str("user_agent", ctx.Request.UserAgent()).
				Str("IP", ctx.RemoteIP()).
				// エラー
				Err(fmt.Errorf("panic: %v", err)).
				Int("ErrCode", ec).
				Str("StackTrace", string(debug.Stack())).
				// パフォーマンス
				Dur("latency(ms)", elapsedTime).
				Msg("")

			// レスポンス作成
			res := struct {
				StatusCode int    `json:"statusCode"`
				Title      string `json:"title"`
				Message    string `json:"message"`
			}{
				StatusCode: http.StatusInternalServerError,
				Title:      "Internal Server Error",
				Message:    "予期せぬエラーが発生しました。\nしばらくしてから再度お試しください。",
			}
			ctx.JSON(http.StatusInternalServerError, res)
		}
	}()
	ctx.Next()
}
