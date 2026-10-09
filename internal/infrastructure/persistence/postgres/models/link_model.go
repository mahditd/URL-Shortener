package models

import "time"

type LinkModel struct {
	ID uint `gorm:"primaryKey"`

	URL  string `gorm:"uniqueIndex;not null"`
	Code string `gorm:"uniqueIndex;not null"`

	CreatedAt time.Time
}

func (LinkModel) TableName() string {
	return "links"
}
