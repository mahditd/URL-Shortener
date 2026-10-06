package ports

import "github.com/mahditd/url-shortener/internal/domain/entities"

type LinkRepository interface {
	Save(link entities.Link) error

	FindByCode(code string) (*entities.Link, error)

	FindByURL(url string) (*entities.Link, error)
}
