package usecase

import (
	"fmt"
	"testing"

	"github.com/mahditd/url-shortener/internal/application/dto"
	"github.com/mahditd/url-shortener/internal/infrastructure/persistence/memory"
)

func BenchmarkShorten(b *testing.B) {

	repo := memory.NewLinkRepository()
	linkUsecase := NewLinkUsecase(repo, "http://localhost:8080")

	b.ResetTimer()

	for i := 0; b.Loop(); i++ {

		req := dto.ShortenRequest{URL: fmt.Sprintf("https://go.dev/doc/%d", i)}

		_, err := linkUsecase.Shorten(req)

		if err != nil {
			b.Fatal(err)
		}

	}
}

func BenchmarkRedirect(b *testing.B) {

	repo := memory.NewLinkRepository()
	linkUsecase := NewLinkUsecase(repo, "http://localhost:8080")

	req := dto.ShortenRequest{URL: "https://go.dev/doc/"}

	res, err := linkUsecase.Shorten(req)

	if err != nil {
		b.Fatal(err)
	}

	savedCode := res.Code

	b.ResetTimer()
	for b.Loop() {

		_, err := linkUsecase.GetLinkByCode(savedCode)

		if err != nil {
			b.Fatal(err)
		}

	}
}
