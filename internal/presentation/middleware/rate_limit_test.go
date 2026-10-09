package middleware

import (
	"testing"
	"time"
)

func TestRateLimiter_AllowsWithinLimit(t *testing.T) {
	limiter := NewRateLimiter(3, time.Minute)

	for i := 0; i < 3; i++ {
		if !limiter.Allow("127.0.0.1") {
			t.Fatal("expected request to be allowed")
		}
	}
}


func TestRateLimiter_BlocksAfterLimit(t *testing.T) {
	limiter := NewRateLimiter(2, time.Minute)

	limiter.Allow("127.0.0.1")
	limiter.Allow("127.0.0.1")

	if limiter.Allow("127.0.0.1") {
		t.Fatal("expected request to be blocked")
	}
}


func TestRateLimiter_ResetsAfterWindow(t *testing.T) {
	limiter := NewRateLimiter(1, 10*time.Millisecond)

	if !limiter.Allow("127.0.0.1") {
		t.Fatal("expected first request to pass")
	}

	if limiter.Allow("127.0.0.1") {
		t.Fatal("expected second request to fail")
	}

	time.Sleep(15 * time.Millisecond)

	if !limiter.Allow("127.0.0.1") {
		t.Fatal("expected request after reset window to pass")
	}
}