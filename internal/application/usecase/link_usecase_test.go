package usecase

import (
	"errors"
	"testing"

	"github.com/mahditd/url-shortener/internal/application/dto"
	domainerrors "github.com/mahditd/url-shortener/internal/domain/errors"
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

func TestGetLinkByCodeSuccess(t *testing.T) {

	repo := memory.NewLinkRepository()

	linkUsecase := NewLinkUsecase(repo, "http://localhost:8080")

	req := dto.ShortenRequest{URL: "https://go.dev/doc/"}

	res, err := linkUsecase.Shorten(req)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	link, err := linkUsecase.GetLinkByCode(res.Code)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if link == nil {
		t.Fatalf("expected link not to be nil")
	}

	if link.URL != "https://go.dev/doc/" {
		t.Fatalf("expected urls to be equal")
	}

	if link.Code != res.Code {
		t.Fatalf("expected codes to be equal")
	}

}

func TestGetLinkByCodeNotFound(t *testing.T) {

	repo := memory.NewLinkRepository()

	linkUsecase := NewLinkUsecase(repo, "http://localhost:8080")

	link, err := linkUsecase.GetLinkByCode("randomThing")

	if !errors.Is(err, domainerrors.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}

	if link != nil {
		t.Fatalf("expected link to be nil")
	}

}

func TestNormalizeURL(t *testing.T) {

	tests := []struct {
		name   string
		input  string
		output string
		err    error
	}{
		{
			name:  "missing url",
			input: "",
			err:   domainerrors.ErrInvalidURL,
		},
		{
			name:  "invalid scheme",
			input: "ftp://example.com",
			err:   domainerrors.ErrInvalidURL,
		},
		{
			name:  "random string",
			input: "hello",
			err:   domainerrors.ErrInvalidURL,
		},
		{
			name:   "normal URL",
			input:  "https://go.dev/doc/",
			output: "https://go.dev/doc/",
		},
		{
			name:   "trim spaces and case insensitivity",
			input:  "   HTTPS://GO.DEV/doc/   ",
			output: "https://go.dev/doc/",
		},
		{
			name:   "default ports : http",
			input:  "HTTP://Example.com:80/path",
			output: "http://example.com/path",
		},
		{
			name:   "default ports : https",
			input:  "HTTPS://Example.com:443/path",
			output: "https://example.com/path",
		},
		{
			name:  "no host",
			input: "http:///path",
			err:   domainerrors.ErrInvalidURL,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := normalizeURL(tt.input)

			if tt.err != nil {
				if err == nil {
					t.Fatalf("expected error %v, got nil", tt.err)
				}

				if !errors.Is(err, tt.err) {
					t.Fatalf("expected %v, got %v", tt.err, err)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}

				if out != tt.output {
					t.Fatalf("expected output %s, got %s", tt.output, out)
				}
			}

		})

	}
}
