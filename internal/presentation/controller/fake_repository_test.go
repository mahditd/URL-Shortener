package controller_test

import (
	"github.com/mahditd/url-shortener/internal/domain/entities"
	"github.com/mahditd/url-shortener/internal/domain/ports"
)

type FakeRepository struct {
	FindByCodeResult *entities.Link
	FindByCodeError  error
	FindByURLResult  *entities.Link
	FindByURLError   error
	SaveError        error
}

func (f *FakeRepository) FindByCode(code string) (*entities.Link, error) {
	return f.FindByCodeResult, f.FindByCodeError
}

func (f *FakeRepository) FindByURL(URL string) (*entities.Link, error) {
	return f.FindByURLResult, f.FindByURLError
}

func (f *FakeRepository) Save(link entities.Link) error {
	return f.SaveError
}

var _ ports.LinkRepository = (*FakeRepository)(nil)
