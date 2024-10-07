package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRepository interface {
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) (string, error)
	Delete(key string) error
}

type redisRepository struct {
	cli *redis.Client
	ctx context.Context
}

func NewRedisClient(cli *redis.Client, ctx context.Context) RedisRepository {
	return &redisRepository{
		cli: cli,
		ctx: ctx,
	}
}

func (r *redisRepository) Set(key string, value interface{}, expiration time.Duration) error {
	return r.cli.Set(r.ctx, key, value, expiration).Err()
}

func (r *redisRepository) Get(key string) (string, error) {
	return r.cli.Get(r.ctx, key).Result()
}

func (r *redisRepository) Delete(key string) error {
	return r.cli.Del(r.ctx, key).Err()
}

// func (r redisRepository) SetSettion(s model.Session) error {
// 	sessionData := map[string]interface{}{
// 		"user_id":    s.UserID,
// 		"created_at": s.CreatedAt,
// 	}
// 	err := r.red.HSet(r.ctx, s.SessionID, sessionData).Err()
// 	if err != nil {
// 		return xerrors.Errorf("redis don't set sessionID: %w", err)
// 	}

// 	// セッションの有効期限を30分に設定
// 	err = r.red.Expire((r.ct, sessionKey, 30*time.Minute).Err()
// 	if err != nil {
// 		return xerrors.Errorf("redis don't set limit time: %w", err)
// 	}

// 	return nil
// }
