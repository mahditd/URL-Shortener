package usecase

import (
	"errors"
	"net/url"
	"time"

	"github.com/mahditd/url-shortener/internal/application/dto"
	"github.com/mahditd/url-shortener/internal/domain/entities"
	domainerrors "github.com/mahditd/url-shortener/internal/domain/errors"
	"github.com/mahditd/url-shortener/internal/domain/ports"
)

type LinkUsecase struct {
	repository ports.LinkRepository
	baseUrl    string
}

func NewLinkUsecase(repository ports.LinkRepository, baseUrl string) *LinkUsecase {
	return &LinkUsecase{
		repository: repository,
		baseUrl:    baseUrl,
	}
}

func (u *LinkUsecase) Shorten(req dto.ShortenRequest) (*dto.ShortenResponse, error) {

	if !validateURL(req.URL) {
		return nil, domainerrors.ErrInvalidURL
	}

	exists, err := u.repository.FindByURL(req.URL)

	if err == nil {
		return &dto.ShortenResponse{
			Code:     exists.Code,
			ShortURL: u.baseUrl + "/" + exists.Code,
		}, nil
	}

	var code string

	for {
		code = generateCode()

		_, err := u.repository.FindByCode(code)

		if errors.Is(err, domainerrors.ErrNotFound) {
			break
		}
	}
	link := entities.Link{Code: code, URL: req.URL, CreatedAt: time.Now()}

	err = u.repository.Save(link)

	if err != nil {
		return nil, err
	}

	return &dto.ShortenResponse{
		Code:     code,
		ShortURL: u.baseUrl + "/" + code,
	}, nil
}

func (u *LinkUsecase) GetLinkByCode(code string) (*entities.Link, error) {

	return u.repository.FindByCode(code)
}

func validateURL(rawURL string) bool {
	if rawURL == "" {
		return false
	}

	parsed, err := url.Parse(rawURL)

	if err != nil {
		return false
	}

	return parsed.Scheme == "http" || parsed.Scheme == "https"
}
