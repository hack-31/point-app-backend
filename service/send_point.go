package service

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/hack-31/point-app-backend/auth"
	"github.com/hack-31/point-app-backend/domain"
	"github.com/hack-31/point-app-backend/domain/model"
	"github.com/hack-31/point-app-backend/repository"
)

type SendPoint struct {
	PointRepo        domain.PointRepo
	UserRepo         domain.UserRepo
	NotificationRepo domain.NotificationRepo
	Connection       *repository.AppConnection
	Cache            domain.Cache
}

func NewSendPoint(repo *repository.Repository, connection *repository.AppConnection, cache domain.Cache) *SendPoint {
	return &SendPoint{PointRepo: repo, UserRepo: repo, NotificationRepo: repo, Connection: connection, Cache: cache}
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
	fromUserID := uid.(model.UserID)

	// トランザクション開始
	if err := sp.Connection.Begin(ctx); err != nil {
		return fmt.Errorf("cannot trasanction: %w ", err)
	}
	defer func() { _ = sp.Connection.Rollback() }()

	// 送付可能か残高を調べる
	email, _ := ctx.Get(auth.Email)
	stringMail := email.(string)
	u, err := sp.UserRepo.FindUserByEmail(ctx, sp.Connection.DB(), &stringMail)
	if err != nil {
		return fmt.Errorf("cannot FindUserByEmail in sending point: %w ", err)
	}
	sendablePoint := model.NewSendablePoint(u.SendingPoint)
	if !sendablePoint.CanSendPoint(sendPoint) {
		return fmt.Errorf("can not send for not having sendable point: %w", repository.ErrHasNotSendablePoint)
	}

	// 送信相手にポイントを加算
	if err := sp.PointRepo.RegisterPointTransaction(ctx, sp.Connection.DB(), fromUserID, model.UserID(toUserId), sendPoint); err != nil {
		return fmt.Errorf("cannot RegisterPointTransaction in sending point: %w ", err)
	}

	// 送信ユーザの送信可能ポイントを減らす
	if err := sp.PointRepo.UpdateSendablePoint(ctx, sp.Connection.DB(), fromUserID, sendablePoint.CalculatePointBalance(sendPoint)); err != nil {
		return fmt.Errorf("cannot updateSendablePoint in sending point: %w ", err)
	}

	// お知らせ内容作成
	fromUser, err := sp.UserRepo.GetUserByID(ctx, sp.Connection.DB(), fromUserID)
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
	registeredNotif, err := sp.NotificationRepo.CreateNotification(ctx, sp.Connection.DB(), n)
	if err != nil {
		return fmt.Errorf("cannot CreateNotification in sending point: %w ", err)
	}

	// 最新お知らせIDをユーザテーブルに保存
	if err := sp.UserRepo.UpdateNotificationLatestIDByID(ctx, sp.Connection.DB(), n.ToUserID, registeredNotif.ID); err != nil {
		return fmt.Errorf("cannot UpdateNotificationLatestIDByID in sending point: %w ", err)
	}

	// トランザクションコミット
	if err := sp.Connection.Commit(); err != nil {
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
