package service

import (
	"encoding/json"
	"fmt"

	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
	"github.com/hack-31/point-app-backend/domain"
	"github.com/hack-31/point-app-backend/domain/model"
	"github.com/hack-31/point-app-backend/myerror"
	"github.com/hack-31/point-app-backend/repository"
	"github.com/hack-31/point-app-backend/repository/entity"
	"github.com/hack-31/point-app-backend/utils"
	"github.com/jmoiron/sqlx"
)

type SendPoint struct {
	PointRepo        domain.PointRepo
	UserRepo         domain.UserRepo
	NotificationRepo domain.NotificationRepo
	Tx               repository.Beginner
	Cache            domain.Cache
}

func NewSendPoint(repo *repository.Repository, db *sqlx.DB, cache domain.Cache) *SendPoint {
	return &SendPoint{PointRepo: repo, UserRepo: repo, NotificationRepo: repo, Tx: db, Cache: cache}
}

// ポイント送信サービス
//
// @params
// ctx コンテキスト
// toUserId 送付先ユーザーID
// sendPoint 送付ポイント
func (sp *SendPoint) SendPoint(ctx *gin.Context, toUserId, sendPoint int) error {
	// コンテキストよりUserIDを取得
	fromUserID := utils.GetUserID(ctx)

	// トランザクション開始
	tx, err := sp.Tx.BeginTxx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}
	// ロールバック
	defer func() { _ = tx.Rollback() }()

	// 送付可能か残高を調べる
	mail := utils.GetEmail(ctx)
	u, err := sp.UserRepo.FindUserByEmail(ctx, tx, mail)
	if err != nil {
		return errors.Wrap(err, "failed to get user by email")
	}
	sendablePoint := model.NewSendablePoint(u.SendingPoint)
	if !sendablePoint.CanSendPoint(sendPoint) {
		return errors.Wrap(myerror.ErrHasNotSendablePoint, "failed to send point")
	}

	// 送信相手にポイントを加算
	if err := sp.PointRepo.RegisterPointTransaction(ctx, tx, fromUserID, entity.UserID(toUserId), sendPoint); err != nil {
		return errors.Wrap(err, "failed to register point transaction")
	}

	// 送信ユーザの送信可能ポイントを減らす
	if err := sp.PointRepo.UpdateSendablePoint(ctx, tx, fromUserID, sendablePoint.CalculatePointBalance(sendPoint)); err != nil {
		return errors.Wrap(err, "failed to update sendable point")
	}

	// お知らせ内容作成
	fromUser, err := sp.UserRepo.GetUserByID(ctx, tx, fromUserID)
	if err != nil {
		return errors.Wrap(err, "failed to get user by id")
	}
	n := entity.Notification{
		TypeID:      model.NotificationTypeSendingPoint,
		ToUserID:    entity.UserID(toUserId),
		FromUserID:  fromUserID,
		Description: fmt.Sprintf("%s%sさんから%dポイント送付されました。", fromUser.FamilyName, fromUser.FirstName, sendPoint),
	}

	// お知らせを新規登録
	registeredNotif, err := sp.NotificationRepo.CreateNotification(ctx, tx, n)
	if err != nil {
		return errors.Wrap(err, "failed to create notification")
	}

	// トランザクションコミット
	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	// お知らせを通知
	channel := fmt.Sprintf("notification:%d", toUserId)
	payload, err := json.Marshal(registeredNotif)
	if err != nil {
		return errors.Wrap(err, "failed to marshal in send point")
	}
	if err := sp.Cache.Publish(ctx, channel, string(payload)); err != nil {
		return errors.Wrapf(err, "failed to publish to %q channel", channel)
	}

	return nil
}
