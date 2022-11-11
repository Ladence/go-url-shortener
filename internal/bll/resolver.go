package bll

import (
	"context"
	"fmt"

	"github.com/Ladence/go-url-shortener/internal/storage"
)

type Resolver struct {
	urlStorage storage.KvStorage
}

func NewResolver(kvStorage storage.KvStorage) *Resolver {
	return &Resolver{
		urlStorage: kvStorage,
	}
}

func (r *Resolver) Resolve(shortUrl string) (string, error) {
	value, err := r.urlStorage.Get(context.Background(), shortUrl)
	if err != nil {
		return "", fmt.Errorf("unable to get key from kvstorage: %v", err)
	}
	return value.(string), nil
}
