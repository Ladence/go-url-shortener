package storage

import (
	"context"
	"testing"
	"time"
)

func TestMap(t *testing.T) {
	mapStorage := NewMapStorage(10)

	cases := []struct {
		name    string
		key     string
		val     string
		keepFor time.Duration
	}{
		{
			"Push and get value",
			"key",
			"value",
			time.Second * 10,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			err := mapStorage.Push(context.Background(), tt.key, tt.val, tt.keepFor)
			if err != nil {
				t.Errorf("failed on push: %v", err)
			}
			got, err := mapStorage.Get(context.Background(), tt.key)
			if err != nil {
				t.Errorf("failed on push: %v", err)
			}
			if got != tt.val {
				t.Errorf("pushed and received values isn't equal, pushed: %s but got: %s", tt.val, got)
			}
		})
	}
}
