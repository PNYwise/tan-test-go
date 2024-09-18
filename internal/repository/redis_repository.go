package repository

import (
	"context"
	"errors"
	"tan-test-go/internal/domain"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	ctx    context.Context
	client *redis.Client
}

func NewRedisRepository(ctx context.Context, client *redis.Client) domain.RedisCacheRepository {
	return &RedisRepository{client: client, ctx: ctx}
}

func (r *RedisRepository) Get(key string) (string, error) {
	return r.client.Get(r.ctx, key).Result()
}

func (r *RedisRepository) Set(key string, value interface{}, expiration time.Duration) error {
	val, ok := value.(string)
	if !ok {
		return errors.New("value must be a string")
	}
	return r.client.Set(r.ctx, key, val, expiration).Err()
}
