package ratelimiter

import (
	"sync"
	"time"
)

type RateLimiter struct {
	maxRequests int
	duration    time.Duration
	requests    sync.Map
}

type requestInfo struct {
	lastRequest time.Time
	counter     int
}

func NewRateLimiter(maxRequests int, duration time.Duration) *RateLimiter {
	return &RateLimiter{
		maxRequests: maxRequests,
		duration:    duration,
	}
}

func (rl *RateLimiter) Allow(identifier string) bool {
	now := time.Now()

	value, _ := rl.requests.LoadOrStore(identifier, requestInfo{lastRequest: now, counter: 0})
	reqInfo := value.(requestInfo)

	if now.Sub(reqInfo.lastRequest) > rl.duration {
		reqInfo.lastRequest = now
		reqInfo.counter = 1
	} else {
		if reqInfo.counter < rl.maxRequests {
			reqInfo.counter++
		} else {
			return false
		}
	}

	rl.requests.Store(identifier, reqInfo)
	return true
}
