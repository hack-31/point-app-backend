package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hack-31/point-app-backend/config"
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

	router := gin.Default()
	// ミドルウェアの設定
	router.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"http://localhost", "https://*.dkjrwfcbom7qp.amplifyapp.com", "https://hack-31.github.io"},
		AllowMethods:  []string{"*"},
		AllowHeaders:  []string{"*"},
		AllowWildcard: true,
	}))

	// DB関係初期化
	db, cleanup, err := repository.New(ctx, cfg)
	if err != nil {
		return err
	}
	defer cleanup()

	err = routers.SetRouting(ctx, db, router, cfg)
	if err != nil {
		return err
	}
	err = routers.SetAuthRouting(ctx, db, router, cfg)
	if err != nil {
		return err
	}

	return router.Run(fmt.Sprintf(":%d", cfg.Port))
}
