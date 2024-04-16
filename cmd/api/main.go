package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/hack-31/point-app-backend/config"
	"github.com/hack-31/point-app-backend/middleware"
	"github.com/hack-31/point-app-backend/repository"
	routers "github.com/hack-31/point-app-backend/router"
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Printf("failed to terminated server: %v", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	cfg, err := config.New()
	if err != nil {
		return err
	}

	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else if cfg.Env == "test" {
		gin.SetMode(gin.TestMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	r := gin.New()
	// ミドルウェアの設定
	r.Use(middleware.Recovery)
	r.Use(middleware.Logger())
	r.Use(middleware.CORSMiddleware())

	// DB関係初期化
	db, cleanup, err := repository.NewDB(ctx, cfg)
	if err != nil {
		return err
	}
	defer cleanup()

	// ルーティング初期化
	if err = routers.SetRouting(ctx, db, r, cfg); err != nil {
		return err
	}

	// サーバー起動
	log.Printf("Listening and serving HTTP on :%v", cfg.Port)
	server := NewServer(r, fmt.Sprintf(":%d", cfg.Port))
	return server.Run(ctx)
}
