package bll

import (
	"context"
	"github.com/Ladence/go-url-shortener/internal/storage"
	"testing"
	"time"
)

func TestResolver(t *testing.T) {
	mStorage := storage.NewMapStorage(10)

	// warm storage with some fake values, just testing only integration
	_ = mStorage.Push(context.Background(), "http://google.com", "TTTy1235", time.Hour*12)

	resolver := NewResolver(mStorage)

	cases := []struct {
		name        string
		inUrl       string
		expectedUrl string
		wantErr     bool
	}{
		{
			"Get existed shorten URL for google.com",
			"http://google.com",
			"TTTy1235",
			false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			shUrl, err := resolver.Resolve(tt.inUrl)
			if shUrl != tt.expectedUrl {
				t.Fatalf("case expectedUrl URL is: %s but got: %s", tt.expectedUrl, shUrl)
			}
			if (err != nil) != tt.wantErr {
				t.Fatalf("case wantErr = %t but err is = %v", tt.wantErr, err)
			}
		})
	}

}
