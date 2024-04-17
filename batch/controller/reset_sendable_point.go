package controller

import (
	"github.com/urfave/cli/v2"
)

type ResetSendablePoint struct {
	Usercase ResetSendablePointer
}

func NewResetSendablePoint(r ResetSendablePointer) *ResetSendablePoint {
	return &ResetSendablePoint{Usercase: r}
}

// ResetSendablePoint は、送付可能なポイントをリセットする
func (rsp *ResetSendablePoint) ResetSendablePoint(cCtx *cli.Context) error {
	initialSendablePoint := cCtx.Int("point")
	if err := rsp.Usercase.ResetPoint(cCtx.Context, initialSendablePoint); err != nil {
		return err
	}
	return nil
}
