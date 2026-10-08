package database

import (
	"github.com/mahditd/url-shortener/internal/infrastructure/persistence/postgres/models"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&models.LinkModel{})
}