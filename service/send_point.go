package service

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/hack-31/point-app-backend/domain"
	"github.com/hack-31/point-app-backend/domain/model"
	"github.com/hack-31/point-app-backend/repository"
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
		return err
	}
	// ロールバック
	defer func() { _ = tx.Rollback() }()

	// 送付可能か残高を調べる
	mail := utils.GetEmail(ctx)
	u, err := sp.UserRepo.FindUserByEmail(ctx, tx, mail)
	if err != nil {
		return fmt.Errorf("cannot FindUserByEmail in sending point: %w ", err)
	}
	sendablePoint := model.NewSendablePoint(u.SendingPoint)
	if !sendablePoint.CanSendPoint(sendPoint) {
		return fmt.Errorf("can not send for not having sendable point: %w", repository.ErrHasNotSendablePoint)
	}

	// 送信相手にポイントを加算
	if err := sp.PointRepo.RegisterPointTransaction(ctx, tx, fromUserID, model.UserID(toUserId), sendPoint); err != nil {
		return fmt.Errorf("cannot RegisterPointTransaction in sending point: %w ", err)
	}

	// 送信ユーザの送信可能ポイントを減らす
	if err := sp.PointRepo.UpdateSendablePoint(ctx, tx, fromUserID, sendablePoint.CalculatePointBalance(sendPoint)); err != nil {
		return fmt.Errorf("cannot updateSendablePoint in sending point: %w ", err)
	}

	// お知らせ内容作成
	fromUser, err := sp.UserRepo.GetUserByID(ctx, tx, fromUserID)
	if err != nil {
		return fmt.Errorf("cannot getUserByID in sending point: %w ", err)
	}
	n := model.Notification{
		TypeID:      model.NotificationTypeSendingPoint,
		ToUserID:    model.UserID(toUserId),
		FromUserID:  fromUserID,
		Description: fmt.Sprintf("%s%sさんから%dポイント送付されました。", fromUser.FamilyName, fromUser.FirstName, sendPoint),
	}

	// お知らせを新規登録
	registeredNotif, err := sp.NotificationRepo.CreateNotification(ctx, tx, n)
	if err != nil {
		return fmt.Errorf("cannot CreateNotification in sending point: %w ", err)
	}

	// トランザクションコミット
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("cannot Comit in send point: %w ", err)
	}

	// お知らせを通知
	channel := fmt.Sprintf("notification:%d", toUserId)
	payload, err := json.Marshal(registeredNotif)
	if err != nil {
		return fmt.Errorf("cannot Marshal in send point :%w", err)
	}
	if err := sp.Cache.Publish(ctx, channel, string(payload)); err != nil {
		return fmt.Errorf("cannot Publish %s channel :%w", channel, err)
	}

	return nil
}
