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
		return nil, res.Err()
	}
	return res.String(), nil
}
