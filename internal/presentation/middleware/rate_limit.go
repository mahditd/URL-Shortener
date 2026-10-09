package middleware

import (
	"sync"
	"time"
)

type client struct {
	requests int
	resetAt  time.Time
}

type RateLimiter struct {
	mu      sync.Mutex
	clients map[string]*client

	limit  int
	window time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		clients: make(map[string]*client),
		limit:   limit,
		window:  window,
	}
}

func (r *RateLimiter) Allow(ip string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()

	c, exists := r.clients[ip]

	if !exists || now.After(c.resetAt) {
		r.clients[ip] = &client{
			requests: 1,
			resetAt:  now.Add(r.window),
		}

		return true
	}

	if c.requests >= r.limit {
		return false
	}

	c.requests++

	return true
}