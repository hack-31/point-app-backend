package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hack-31/point-app-backend/utils/clock"
	"github.com/rs/zerolog"
	"golang.org/x/xerrors"
)

// コンテキストに保存するエラー情報のキー定義
const (
	errKey     = "err"
	errCodeKey = "errCodo"
	rpsKey     = "rps"
)

type Log struct {
	errCode int
	rps     []byte
	err     error
}

// ログに保存する構造体の作成
// 最後に登録するにはLoggingを呼ぶ必要あり
func NewLog() *Log {
	return &Log{}
}
func (l *Log) ErrCode(code int) *Log {
	l.errCode = code
	return l
}
func (l *Log) Rsp(r []byte) *Log {
	l.rps = r
	return l
}
func (l *Log) Err(e error) *Log {
	l.err = e
	return l
}

// ログ情報を登録
func (l *Log) Logging(ctx *gin.Context) {
	ctx.Set(rpsKey, l.rps)
	ctx.Set(errCodeKey, l.errCode)
	ctx.Set(errKey, l.err)
}

// Logger は、ロギング処理
func Logger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set("startTime", time.Now())
		path := ctx.Request.URL.Path
		raw := ctx.Request.URL.RawQuery
		if raw != "" {
			path = path + "?" + raw
		}

		ctx.Next()

		// エラーデータ取得
		status := ctx.Writer.Status()

		gec, exist := ctx.Get(errCodeKey)
		var ec int
		if exist {
			ec = gec.(int)
		}
		gr, exists := ctx.Get(rpsKey)
		var rps []byte
		if exists {
			rps = gr.([]byte)
		}

		ge, exists := ctx.Get(errKey)
		var e error
		if exists && ge != nil {
			e = ge.(error)
		}

		gst, exists := ctx.Get("startTime")
		var startTime time.Time
		if exists {
			startTime = gst.(time.Time)
		}
		elapsedTime := time.Since(startTime)

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

		if status >= 400 && status < 500 {
			logger.Warn().
				// レスポンス
				Int("status", status).
				Bytes("response", rps).
				// リクエスト
				Str("URI", path).
				Str("method", ctx.Request.Method).
				Str("user_agent", ctx.Request.UserAgent()).
				Str("IP", ctx.RemoteIP()).
				// エラー
				Err(e).
				Int("ErrCode", ec).
				Str("StackTrace", fmt.Sprintf("%+v", xerrors.Errorf("%w", e))).
				// パフォーマンス
				Dur("latency(ms)", elapsedTime).
				Msg("")
			return
		}
		if status >= 500 {
			logger.Error().
				// レスポンス
				Int("status", status).
				Bytes("response", rps).
				// リクエスト
				Str("URI", path).
				Str("method", ctx.Request.Method).
				Str("user_agent", ctx.Request.UserAgent()).
				Str("IP", ctx.RemoteIP()).
				// エラー
				Err(e).
				Int("ErrCode", ec).
				Str("StackTrace", fmt.Sprintf("%+v", xerrors.Errorf("%w", e))).
				// パフォーマンス
				Dur("latency(ms)", elapsedTime).
				Msg("")
			return
		}

		// 成功時
		logger.Info().
			// レスポンス
			Int("status", status).
			Bytes("response", rps).
			// リクエスト
			// Str("URI", ctx.Request.URL.String()).
			Str("URI", path).
			Str("method", ctx.Request.Method).
			Str("user_agent", ctx.Request.UserAgent()).
			Str("IP", ctx.RemoteIP()).
			// パフォーマンス
			// TODO: runtime, runtime/pprofパッケージを利用してメモリ利用率など出す
			Dur("latency(ms)", elapsedTime).
			Msg("")
	}
}
