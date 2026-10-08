package usecase

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/mahditd/url-shortener/internal/application/dto"
	"github.com/mahditd/url-shortener/internal/domain/entities"
	domainerrors "github.com/mahditd/url-shortener/internal/domain/errors"
	"github.com/mahditd/url-shortener/internal/domain/ports"
)

type LinkUsecase struct {
	repository ports.LinkRepository
	baseUrl    string
	mu         sync.Mutex
}

func NewLinkUsecase(repository ports.LinkRepository, baseUrl string) *LinkUsecase {
	return &LinkUsecase{
		repository: repository,
		baseUrl:    baseUrl,
	}
}

func (u *LinkUsecase) Shorten(req dto.ShortenRequest) (*dto.ShortenResponse, error) {

	normalizedURL, err := normalizeURL(req.URL)

	if err != nil {
		return nil, fmt.Errorf("validate url: %w", err)
	}

	u.mu.Lock()
	defer u.mu.Unlock()

	exists, err := u.repository.FindByURL(normalizedURL)

	if err == nil {
		return &dto.ShortenResponse{
			Code:     exists.Code,
			ShortURL: u.baseUrl + "/" + exists.Code,
		}, nil
	}

	if !errors.Is(err, domainerrors.ErrNotFound) {
		return nil, fmt.Errorf("find existing url: %w", err)
	}

	var code string

	for {
		code = generateCode()

		_, err := u.repository.FindByCode(code)

		if errors.Is(err, domainerrors.ErrNotFound) {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("check code collision: %w", err)
		}
	}
	link := entities.Link{Code: code, URL: normalizedURL, CreatedAt: time.Now()}

	err = u.repository.Save(link)

	if err != nil {
		return nil, fmt.Errorf("save link: %w", err)
	}

	return &dto.ShortenResponse{
		Code:     code,
		ShortURL: u.baseUrl + "/" + code,
	}, nil
}

func (u *LinkUsecase) GetLinkByCode(code string) (*entities.Link, error) {

	link, err := u.repository.FindByCode(code)

	if err != nil {
		return nil, fmt.Errorf("find link: %w", err)
	}

	return link, nil
}

func normalizeURL(rawURL string) (string, error) {

	rawURL = strings.TrimSpace(rawURL)

	if rawURL == "" {
		return "", domainerrors.ErrInvalidURL
	}

	parsed, err := url.Parse(rawURL)

	if err != nil {
		return "", domainerrors.ErrInvalidURL
	}

	parsed.Scheme = strings.ToLower(parsed.Scheme)
	parsed.Host = strings.ToLower(parsed.Host)

	if !(parsed.Scheme == "http" || parsed.Scheme == "https") {
		return "", domainerrors.ErrInvalidURL
	}
	if parsed.Host == "" {
		return "", domainerrors.ErrInvalidURL
	}

	return parsed.String(), nil

}
