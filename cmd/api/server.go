package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

type Server struct {
	srv *http.Server
}

func NewServer(mux http.Handler, addr string) *Server {
	return &Server{
		srv: &http.Server{Handler: mux, Addr: addr},
	}
}

// サーバー起動
//
// @params
// ctx コンテキスト
func (s *Server) Run(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		// http.ErrServerClosed は
		// http.Server.Shutdown() が正常に終了したことを示すので異常ではない
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("failed to close: %+v", err)
			return err
		}
		// mispel := "hoge"
		// fmt.Println(mispel)
		return nil
	})

	// チャンネルからの終了通知を待つ
	<-ctx.Done()
	// サーバーをシャットダウン
	if err := s.srv.Shutdown(context.Background()); err != nil {
		log.Printf("failed to shutdown: %+v", err)
	}
	// 別ルーチンのグレースフルシャットダウンの終了を待つ
	return eg.Wait()
}
