package storage

import (
	"context"
	"math/rand"
	"testing"
	"time"
)

// TestMap_Push tests pushed and extracted value are the same
func TestMap_Push(t *testing.T) {
	mapStorage := NewMapStorage(10)

	cases := []struct {
		name    string
		key     string
		val     any
		keepFor time.Duration
	}{
		{
			"Push and get value (string)",
			"key",
			"value",
			time.Second * 10,
		},
		{
			"Push and get value (numeric)",
			"key",
			15,
			time.Second * 5,
		},
		{
			"Push and get value (arbitrary struct obj)",
			"key12452352",
			struct {
				foo string
				bar int
			}{
				"foo",
				rand.Int(),
			},
			time.Second * 1,
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
