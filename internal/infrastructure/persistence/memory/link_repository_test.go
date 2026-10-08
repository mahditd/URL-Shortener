package memory

import (
	"errors"
	"testing"
	"time"

	"github.com/mahditd/url-shortener/internal/domain/entities"
	domainerrors "github.com/mahditd/url-shortener/internal/domain/errors"
)

func TestSaveAndFindLink(t *testing.T) {
	linkRepo := NewLinkRepository()

	link := entities.Link{URL: "https://go.dev/doc/", Code: "abc123", CreatedAt: time.Now()}

	err := linkRepo.Save(link)

	if err != nil {
		t.Fatalf("expected no errors saving")
	}

	linkByURL, err := linkRepo.FindByURL("https://go.dev/doc/")

	if err != nil {
		t.Fatalf("expected no errors finding by URL")
	}

	linkByCode, err := linkRepo.FindByCode(linkByURL.Code)

	if err != nil {
		t.Fatalf("expected no errors finding by code")
	}

	if linkByURL.URL != link.URL {
		t.Fatalf("expected link's URLs to be the same")
	}

	if linkByURL.Code != link.Code {
		t.Fatalf("expected link's codes to be the same")
	}

	if linkByURL.CreatedAt != link.CreatedAt {
		t.Fatalf("expected links created time to be the same")
	}

	if link.Code != linkByCode.Code {
		t.Fatalf("expected link's codes to be the same")
	}

	if link.URL != linkByCode.URL {
		t.Fatalf("expected link's URLs to be the same")
	}

	if link.CreatedAt != linkByCode.CreatedAt {
		t.Fatalf("expected links created time to be the same")
	}

}

func TestFindByCodeNotFound(t *testing.T) {

	linkRepo := NewLinkRepository()

	link, err := linkRepo.FindByCode("random")

	if err != nil {
		if !errors.Is(err, domainerrors.ErrNotFound) {
			t.Fatalf("FindByCode : expected %v got %v", domainerrors.ErrNotFound, err)
		}
	} else {
		t.Fatalf("FindByCode : expected %v got no error", domainerrors.ErrNotFound)
	}

	if link != nil {
		t.Fatalf("FindByCode : expected link to be nil")
	}

}

func TestFindByURLNotFound(t *testing.T) {

	linkRepo := NewLinkRepository()

	link, err := linkRepo.FindByURL("random")

	if err != nil {
		if !errors.Is(err, domainerrors.ErrNotFound) {
			t.Fatalf("FindByURL : expected %v got %v", domainerrors.ErrNotFound, err)
		}
	} else {
		t.Fatalf("FindByURL : expected %v got no error", domainerrors.ErrNotFound)
	}

	if link != nil {
		t.Fatalf("FindByURL : expected link to be nil")
	}

}
