package bll

import (
	"context"
	"fmt"
	"github.com/Ladence/go-url-shortener/internal/model"
	"github.com/Ladence/go-url-shortener/internal/storage"
	"math/rand"
)

type Shortener struct {
	urlStorage storage.KvStorage
}

func NewShortener(kvStorage storage.KvStorage) (*Shortener, error) {
	return &Shortener{
		urlStorage: kvStorage,
	}, nil
}

func (s *Shortener) ShortenUrl(request *model.GetShortenRequest) (*model.GetShortenResponse, error) {
	id := request.CustomShort
	if len(request.CustomShort) == 0 {
		id = EncodeBase62(rand.Uint64())
	}

	val, err := s.urlStorage.Get(context.Background(), id)
	if err != nil {
		return nil, fmt.Errorf("error on get value by id from kvstorage: %v", err)
	}
	if val != nil {
		return nil, fmt.Errorf("URL custom short is already in use")
	}
	if request.Expiry == 0 {
		request.Expiry = 24
	}

	err = s.urlStorage.Push(context.Background(), id, request.Url, request.Expiry)
	if err != nil {
		return nil, fmt.Errorf("error on push shortened url in kvstorage: %v", err)
	}

	return &model.GetShortenResponse{
		Url:             request.Url,
		CustomShort:     "",
		Expiry:          request.Expiry,
		XRateRemaining:  0,
		XRateLimitReset: 30,
	}, nil
}
