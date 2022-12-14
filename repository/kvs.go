package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/hack-31/point-app-backend/config"
)

type CacheType int

const (
	JWT                   CacheType = 0
	TemporaryUserRegister CacheType = 1
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
