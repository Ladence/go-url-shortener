package model

import "time"

type PostShortenRequest struct {
	Url         string         `json:"url"`
	CustomShort string         `json:"short"`
	Expiry      *time.Duration `json:"expiry,omitempty"`
}

type PostShortenResponse struct {
	Url             string        `json:"url"`
	CustomShort     string        `json:"customShort"`
	Expiry          time.Duration `json:"expiry"`
	XRateLimitReset time.Duration `json:"rate_limit_reset"`
	XRateRemaining  int64         `json:"rate_limit"`
}
