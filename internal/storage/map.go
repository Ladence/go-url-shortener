package storage

import (
	"context"
	"sync"
	"time"
)

type value struct {
	val any
	ttl time.Duration
}

// Map is a key-value in memory storage just in map[string]string struct
type Map struct {
	mtx  sync.RWMutex
	m    map[string]value
	stop chan struct{}
}

func NewMapStorage(initialCapacity int) *Map {
	return &Map{
		m: make(map[string]value, initialCapacity),
	}
}

//func (m *Map) RunCleanup() {
//	go func() {
//		select {
//		case <-m.stop:
//			return
//		}
//	}()
//}

func (m *Map) Push(ctx context.Context, key string, val any, duration time.Duration) error {
	m.mtx.Lock()
	m.m[key] = value{val: val, ttl: duration}
	m.mtx.Unlock()
	return nil
}

func (m *Map) Get(ctx context.Context, key string) (any, error) {
	v := m.getInternalValue(key)
	if v == nil {
		return nil, nil
	}
	return v.val, nil
}

func (m *Map) Incr(ctx context.Context, key string) (any, error) {
	panic("not supported yet")
}

func (m *Map) Decr(ctx context.Context, key string) (any, error) {
	panic("not supported yet")
}

func (m *Map) TTL(ctx context.Context, key string) (time.Duration, error) {
	v := m.getInternalValue(key)
	if v == nil {
		return 0, nil
	}
	return v.ttl, nil
}

func (m *Map) Close() error {
	for k := range m.m {
		delete(m.m, k)
	}
	return nil
}

func (m *Map) getInternalValue(key string) *value {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	val, ok := m.m[key]
	if !ok {
		return nil
	}
	return &val
}
