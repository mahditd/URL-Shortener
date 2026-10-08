package usecase

import (
	"testing"

	"github.com/mahditd/url-shortener/internal/application/dto"
	"github.com/mahditd/url-shortener/internal/infrastructure/persistence/memory"
)

func TestShortenSuccess(t *testing.T) {

	repo := memory.NewLinkRepository()

	linkUsecase := NewLinkUsecase(repo, "http://localhost:8080")

	req := dto.ShortenRequest{URL: "https://go.dev/doc/"}

	res, err := linkUsecase.Shorten(req)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res == nil {
		t.Fatal("expected response, got nil")
	}

	if res.Code == "" {
		t.Fatal("expected code, got empty")
	}

	if res.ShortURL != "http://localhost:8080/"+res.Code {
		t.Fatalf("unexpected short url: %s", res.ShortURL)
	}

}

func TestShortenIdempotency(t *testing.T) {

	repo := memory.NewLinkRepository()

	linkUsecase := NewLinkUsecase(repo, "http://localhost:8080")

	req := dto.ShortenRequest{URL: "https://go.dev/doc/"}

	res1, err := linkUsecase.Shorten(req)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	res2, err := linkUsecase.Shorten(req)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res1.Code != res2.Code {
		t.Fatalf("expected equal codes")
	}

	if res1.ShortURL != res2.ShortURL {
		t.Fatalf("expected equal short urls")
	}

}
