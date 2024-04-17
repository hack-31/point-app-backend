package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hack-31/point-app-backend/batch/wire"
	"github.com/hack-31/point-app-backend/utils/clock"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
	"golang.org/x/xerrors"
)

func main() {
	app := &cli.App{
		Name:  "batch",
		Usage: "point app batch processing",
		Commands: []*cli.Command{
			{
				Name:        "reset_sendable_point",
				Description: "[月初のの24:00] 送付可能なポイントをリセットする",
				Usage:       "Reset sendable point",
				UsageText:   "batch reset_sendable_point [args]",
				Aliases:     []string{"rsp"},
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:        "point",
						Aliases:     []string{"p"},
						Usage:       "Set default sendable point",
						Value:       10000,
						DefaultText: "10000",
					},
				},
				Action: func(cCtx *cli.Context) error {
					s, cleanup, err := wire.InitResetPointablePoint(cCtx.Context)
					defer cleanup()
					if err != nil {
						logErr(cCtx, err)
						return err
					}
					if err := s.ResetSendablePoint(cCtx); err != nil {
						logErr(cCtx, err)
						return err
					}
					logSuccess(cCtx)
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func logErr(cCtx *cli.Context, err error) {
	option := ""
	for _, v := range cCtx.FlagNames() {
		option += fmt.Sprintf("%s=%s ", v, cCtx.String(v))
	}

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
		// リクエスト
		Str("command", cCtx.Command.Name).
		Str("options", option).
		// エラー
		Err(err).
		Str("StackTrace", fmt.Sprintf("%+v", xerrors.Errorf("%w", err))).
		Msg("")
}

func logSuccess(cCtx *cli.Context) {
	// オプション
	option := ""
	for _, v := range cCtx.FlagNames() {
		option += fmt.Sprintf("%s=%s ", v, cCtx.String(v))
	}

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
	logger.Info().
		// リクエスト
		Str("command", cCtx.Command.Name).
		Str("options", option).
		Msg("")
}
