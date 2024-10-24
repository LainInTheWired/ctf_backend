package repository

import (
	"context"
	"time"

	"github.com/cockroachdb/errors"
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
	err := r.cli.Set(r.ctx, key, value, expiration).Err()
	if err != nil {
		return errors.Wrap(err, "redis can't set: %w")
	}
	return nil
}

func (r *redisRepository) Get(key string) (string, error) {
	s, err := r.cli.Get(r.ctx, key).Result()
	if err != nil {
		return "", errors.Wrap(err, "error set redis")
	}
	return s, nil
}

func (r *redisRepository) Delete(key string) error {
	return errors.Wrap(r.cli.Del(r.ctx, key).Err(), "error delete redis")
}
