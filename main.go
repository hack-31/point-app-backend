package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hack-31/point-app-backend/config"
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
		// TODO: 一旦全てのドメインを許可
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"*"},
		AllowHeaders:  []string{"*"},
		AllowWildcard: true,
	}))

	routers.SetRouting(ctx, router)
	return router.Run(fmt.Sprintf(":%d", cfg.Port))
}
