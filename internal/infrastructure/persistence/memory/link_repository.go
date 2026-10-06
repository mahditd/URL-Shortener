package memory

import (
	"sync"

	"github.com/mahditd/url-shortener/internal/domain/entities"
	domainerrors "github.com/mahditd/url-shortener/internal/domain/errors"
	"github.com/mahditd/url-shortener/internal/domain/ports"
)

type LinkRepository struct {
	mu    sync.RWMutex
	links map[string]entities.Link
	urls  map[string]string
}

func NewLinkRepository() ports.LinkRepository {
	return &LinkRepository{
		links: make(map[string]entities.Link),
		urls:  make(map[string]string),
	}
}

func (r *LinkRepository) Save(link entities.Link) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.links[link.Code] = link
	r.urls[link.URL] = link.Code

	return nil
}

func (r *LinkRepository) FindByCode(code string) (*entities.Link, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	link, exists := r.links[code]

	if !exists {
		return nil, domainerrors.ErrNotFound
	}

	return &link, nil
}

func (r *LinkRepository) FindByURL(url string) (*entities.Link, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	code, exists := r.urls[url]

	if !exists {
		return nil, domainerrors.ErrNotFound
	}

	link := r.links[code]

	return &link, nil
}
