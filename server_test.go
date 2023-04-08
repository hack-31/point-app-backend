package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"golang.org/x/sync/errgroup"
)

func TestServer_Run(t *testing.T) {
	// 返り値を設定
	ctx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)
	mux := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
	})
	const port = 14280

	// テスト対象関数を呼び出し
	eg.Go(func() error {
		s := NewServer(mux, fmt.Sprintf(":%d", port))
		return s.Run(ctx)
	})

	// HACK: サーバーが起動してから、リクエストを送らないたとエラーになるた50ms待つ
	// 一時的な処理
	timer1 := time.NewTimer(50 * time.Millisecond)
	<-timer1.C

	// GETAPIをリクエスト
	in := "healthcheck"
	url := fmt.Sprintf("http://localhost:%d/%s", port, in)
	t.Logf("try request to %q", url)
	rsp, err := http.Get(url)
	if err != nil {
		t.Errorf("failed to get: %+v", err)
	}

	defer rsp.Body.Close()
	// レスポンス整形
	got, err := io.ReadAll(rsp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}
	// サーバの終了動作を検証する
	cancel()
	if err := eg.Wait(); err != nil {
		t.Fatal(err)
	}

	// 戻り値を検証する
	want := fmt.Sprintf("Hello, %s!", in)
	if string(got) != want {
		t.Errorf("want %q, but got %q", want, got)
	}
}
