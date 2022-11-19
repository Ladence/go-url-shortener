package bll

import (
	"github.com/Ladence/go-url-shortener/internal/storage"
	"testing"
	"time"
)

func TestShortener_ShortenUrl(t *testing.T) {
	shortener := NewShortener(storage.NewMapStorage(10))

	cases := []struct {
		name        string
		inputUrl    string
		customUrl   string
		expiry      time.Duration
		expectedUrl string
		wantErr     bool
	}{
		{
			"Provide not a URL",
			"foobarbaz",
			"",
			time.Hour * 12,
			"",
			true,
		},
		{
			"Provide Custom URL",
			"http://google.com",
			"tttasssdagoo",
			time.Hour,
			"tttasssdagoo",
			false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			shortened, err := shortener.ShortenUrl(tt.inputUrl, tt.customUrl, tt.expiry)
			if shortened != shortened {
				t.Fatalf("case expected URL is: %s but got: %s", tt.expectedUrl, shortened)
			}
			if (err != nil) != tt.wantErr {
				t.Fatalf("case wantErr = %t but err is = %v", tt.wantErr, err)
			}
		})
	}
}
