package postgres

import (
	"github.com/mahditd/url-shortener/internal/domain/entities"
	"github.com/mahditd/url-shortener/internal/infrastructure/persistence/postgres/models"
	"gorm.io/gorm"
)

type LinkRepository struct {
	db *gorm.DB
}

func NewLinkRepository(db *gorm.DB) *LinkRepository {
	return &LinkRepository{
		db: db,
	}
}

func (r *LinkRepository) Save(link entities.Link) error {
	model := toModel(link)

	return r.db.Create(&model).Error
}

func (r *LinkRepository) FindByCode(code string) (*entities.Link, error) {
	var model models.LinkModel

	err := r.db.
		Where("code = ?", code).
		First(&model).
		Error

	if err != nil {
		return nil, err
	}

	link := toEntity(model)

	return &link, nil
}

func (r *LinkRepository) FindByURL(url string) (*entities.Link, error) {
	var model models.LinkModel

	err := r.db.
		Where("url = ?", url).
		First(&model).
		Error

	if err != nil {
		return nil, err
	}

	link := toEntity(model)

	return &link, nil
}

func toModel(link entities.Link) models.LinkModel {
	return models.LinkModel{
		Code:      link.Code,
		URL:       link.URL,
		CreatedAt: link.CreatedAt,
	}
}

func toEntity(model models.LinkModel) entities.Link {
	return entities.Link{
		Code:      model.Code,
		URL:       model.URL,
		CreatedAt: model.CreatedAt,
	}
}