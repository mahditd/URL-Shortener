package controller

import (
	"github.com/mahditd/url-shortener/internal/application/dto"
	"github.com/mahditd/url-shortener/internal/domain/entities"
)

type LinkService interface {
	Shorten(req dto.ShortenRequest) (*dto.ShortenResponse, error)
	GetLinkByCode(code string) (*entities.Link, error)
}
