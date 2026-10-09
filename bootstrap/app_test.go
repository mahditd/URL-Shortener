package bootstrap_test

import (
	"testing"

	"github.com/mahditd/url-shortener/bootstrap"
	"github.com/mahditd/url-shortener/internal/infrastructure/persistence/memory"
	"github.com/mahditd/url-shortener/internal/presentation/middleware"
)

func TestNewApp(t *testing.T) {
	repo := memory.NewLinkRepository()

	rateLimiter := middleware.NewRateLimiter(10, 1)

	app := bootstrap.NewApp(
		"http://localhost:8080",
		repo,
		rateLimiter,
	)

	if app == nil {
		t.Fatal("expected gin engine, got nil")
	}
}