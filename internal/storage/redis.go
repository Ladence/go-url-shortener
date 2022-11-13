package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisStorage struct {
	redisClient *redis.Client
}

func (r *RedisStorage) TTL(ctx context.Context, key string) (time.Duration, error) {
	return r.redisClient.TTL(ctx, key).Result()
}

func (r *RedisStorage) Incr(ctx context.Context, key string) (any, error) {
	return r.redisClient.Incr(ctx, key).Result()
}

func (r *RedisStorage) Decr(ctx context.Context, key string) (any, error) {
	return r.redisClient.Decr(ctx, key).Result()
}

func NewRedisStorage(address string, password string, dbNumber int) (KvStorage, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       dbNumber,
	})

	cmdRes := rdb.Ping(context.Background())
	if err := cmdRes.Err(); err != nil {
		return nil, fmt.Errorf("error on creation redis v8 client: %v", err)
	}
	return &RedisStorage{redisClient: rdb}, nil
}

func (r *RedisStorage) Push(ctx context.Context, key string, value any, duration time.Duration) error {
	return r.redisClient.Set(ctx, key, value, duration).Err()
}

// Get returns only string value for specific key
func (r *RedisStorage) Get(ctx context.Context, key string) (any, error) {
	res := r.redisClient.Get(ctx, key)
	if res.Err() != nil {
		if res.Err() == redis.Nil {
			return nil, nil
		}
		return nil, res.Err()
	}
	return res.Val(), nil
}

func (r *RedisStorage) Close() error {
	return r.redisClient.Close()
}
