package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"github.com/hack-31/point-app-backend/config"
)

type CacheType int

const (
	JWT               CacheType = 0
	TemporaryRegister CacheType = 1
)

func NewKVS(ctx context.Context, cfg *config.Config, t CacheType) (*KVS, error) {
	cli := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort),
		DB:   int(t),
	})
	if err := cli.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return &KVS{Cli: cli}, nil
}

type KVS struct {
	Cli *redis.Client
}

// 値を保存する
//
// @params
// ctx context
// key key
// value value
// minute 有効期限(分)
func (k *KVS) Save(ctx context.Context, key, value string, minute time.Duration) error {
	return k.Cli.Set(ctx, key, value, minute*time.Minute).Err()
}

// 値をロードする
//
// @params
// ctx context
// key key
// value
func (k *KVS) Load(ctx context.Context, key string) (string, error) {
	value, err := k.Cli.Get(ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("failed to get by %q: %w", key, ErrNotFoundSession)
	}
	return value, nil
}

// 値を削除
//
// @params
// ctx context
// key key
func (k *KVS) Delete(ctx context.Context, key string) error {
	_, err := k.Cli.Del(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("failed to delete by %q: %w", key, ErrNotFoundSession)
	}
	return nil
}

// 有効期限を延長
func (k *KVS) Expire(ctx context.Context, key string, minitue time.Duration) error {
	_, err := k.Cli.Expire(ctx, key, minitue*time.Minute).Result()
	if err != nil {
		return fmt.Errorf("failed to expire by %q: %w", key, ErrNotFoundSession)
	}
	return nil
}

// チャネルにパブリッシュする
// @params
// ctx context
// channel チャンネル名
// payload 送信するデータ
func (k *KVS) Publish(ctx context.Context, channel, palyload string) error {
	if err := k.Cli.Publish(ctx, channel, palyload).Err(); err != nil {
		return fmt.Errorf("failed to publish by %q: %w", channel, err)
	}
	return nil
}

// チャネルをサブスクライブする
// @params
// ctx context
// channel チャンネル名
//
// @return
// payload 送信されたら発火し、payloadが送られる
func (k *KVS) Subscribe(ctx *gin.Context, channel string) (<-chan string, error) {
	ch := k.Cli.Subscribe(ctx, channel).Channel()
	payload := make(chan string)

	go func() {
		defer close(payload)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ctx.Request.Context().Done():
				return
			case c, ok := <-ch:
				// ch チャンネルがクローズ
				if !ok {
					return
				}
				payload <- c.Payload
			}
		}
	}()

	return payload, nil
}
