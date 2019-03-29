package ratelimit

import "net/http"

type RateLimit struct{}

func NewRateLimit() *RateLimit {
	return &RateLimit{}
}

func (rl *RateLimit) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// todo: rate limit
	next(rw, r)
}
