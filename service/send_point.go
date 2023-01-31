package service

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/hack-31/point-app-backend/auth"
	"github.com/hack-31/point-app-backend/domain"
	"github.com/hack-31/point-app-backend/domain/user"
	"github.com/hack-31/point-app-backend/entity"
	"github.com/hack-31/point-app-backend/repository"
)

type SendPoint struct {
	PointRepo  domain.PointRepo
	UserRepo   domain.UserRepo
	Connection *repository.AppConnection
	DB         repository.Queryer
}

// ポイント送信サービス
//
// @params
// ctx コンテキスト
// toUserId 送付先ユーザーID
// sendPoint 送付ポイント
func (sp *SendPoint) SendPoint(ctx *gin.Context, toUserId, sendPoint int) error {
	// コンテキストよりUserIDを取得
	uid, _ := ctx.Get(auth.UserID)
	fromUserID := uid.(entity.UserID)

	// トランザクション開始
	if err := sp.Connection.Begin(ctx); err != nil {
		return fmt.Errorf("cannot trasanction: %w ", err)
	}

	// 送付可能か残高を調べる
	email, _ := ctx.Get(auth.Email)
	stringMail := email.(string)
	u, err := sp.UserRepo.FindUserByEmail(ctx, sp.DB, &stringMail)
	if err != nil {
		if err := sp.Connection.Rollback(); err != nil {
			return fmt.Errorf("cannot trasanction: %w ", err)
		}
		return err
	}
	sendablePoint := user.NewSendablePoint(u.SendingPoint)
	if !sendablePoint.CanSendPoint(sendPoint) {
		if err := sp.Connection.Rollback(); err != nil {
			return fmt.Errorf("cannot trasanction: %w ", err)
		}
		return fmt.Errorf("can not send for not having sendable point: %w", repository.ErrHasNotSendablePoint)
	}

	// ポイント登録
	if err := sp.PointRepo.RegisterPointTransaction(ctx, sp.Connection.Tx, fromUserID, entity.UserID(toUserId), sendPoint); err != nil {
		if err := sp.Connection.Rollback(); err != nil {
			return fmt.Errorf("cannot trasanction: %w ", err)
		}
		return err
	}

	// 送信ユーザの送信可能ポイントを減らす
	if err := sp.PointRepo.UpdateSendablePoint(ctx, sp.Connection.Tx, fromUserID, sendablePoint.CalculatePointBalance(sendPoint)); err != nil {
		if err := sp.Connection.Rollback(); err != nil {
			return fmt.Errorf("cannot trasanction: %w ", err)
		}
		return err
	}

	// トランザクションコミット
	if err := sp.Connection.Commit(); err != nil {
		return fmt.Errorf("cannot trasanction: %w ", err)
	}

	return nil
}
