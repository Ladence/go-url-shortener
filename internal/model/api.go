package model

import "time"

type GetShortenRequest struct {
	Url         string        `json:"url"`
	CustomShort string        `json:"short"`
	Expiry      time.Duration `json:"expiry"`
}

type GetShortenResponse struct {
	Url             string        `json:"url"`
	CustomShort     string        `json:"customShort"`
	Expiry          time.Duration `json:"expiry"`
	XRateLimitReset time.Duration `json:"rate_limit_reset"`
	XRateRemaining  int           `json:"rate_limit"`
}
