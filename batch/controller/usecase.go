package controller

import (
	"context"
)

type ResetSendablePointer interface {
	ResetPoint(ctx context.Context, initialSendablePoint int) error
}
