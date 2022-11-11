package bll

import (
	"context"
	"fmt"
	"github.com/Ladence/go-url-shortener/internal/storage"
	"github.com/asaskevich/govalidator"
	"math/rand"
	"time"
)

type Shortener struct {
	urlStorage storage.KvStorage
}

func NewShortener(kvStorage storage.KvStorage) (*Shortener, error) {
	return &Shortener{
		urlStorage: kvStorage,
	}, nil
}

// ShortenUrl process incoming url, writes to a key-value storage (for cache) and return a new one shorten url (if custom-not-provided)
func (s *Shortener) ShortenUrl(url, customShortUrl string, expiry *time.Duration) (string, error) {
	if !govalidator.IsURL(url) {
		return "", fmt.Errorf("provided url in request is not a url")
	}

	id := customShortUrl
	if len(id) == 0 {
		id = EncodeBase62(rand.Uint64())
	}
	val, err := s.urlStorage.Get(context.Background(), id)
	if err != nil {
		return "", fmt.Errorf("error on get value by id from kvstorage: %v", err)
	}
	if val != nil {
		return "", fmt.Errorf("URL custom short is already in use")
	}

	expiryStorage := time.Hour * 24
	if expiry != nil {
		expiryStorage = *expiry
	}

	err = s.urlStorage.Push(context.Background(), id, url, expiryStorage)
	if err != nil {
		return "", fmt.Errorf("error on push shortened url in kvstorage: %v", err)
	}

	return id, nil
}
