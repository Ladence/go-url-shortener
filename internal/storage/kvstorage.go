package storage

import (
	"context"
	"time"
)

type KvStorage interface {
	Push(ctx context.Context, key string, value any, duration time.Duration) error
	Get(ctx context.Context, key string) (any, error)
}
