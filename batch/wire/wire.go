//go:build wireinject
// +build wireinject

package wire

import (
	"context"

	"github.com/google/wire"
	"github.com/hack-31/point-app-backend/batch/controller"
	"github.com/hack-31/point-app-backend/batch/usecase"
	"github.com/hack-31/point-app-backend/config"
	"github.com/hack-31/point-app-backend/repository"
	"github.com/hack-31/point-app-backend/utils/clock"
)

func InitResetPointablePoint(ctx context.Context) (*controller.ResetSendablePoint, func(), error) {
	wire.Build(
		config.New,
		repository.NewDB,
		clock.NewRealClocker,
		repository.NewRepository,
		controller.NewResetSendablePoint,
		usecase.NewResetSendablePoint,
		wire.Bind(new(controller.ResetSendablePointer), new(*usecase.ResetSendablePoint)),
	)
	return &controller.ResetSendablePoint{}, nil, nil
}
