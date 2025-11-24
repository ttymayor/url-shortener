package model

import (
	"time"

	"gorm.io/gorm"
)

type URL struct {
	gorm.Model
	OriginalURL string `gorm:"type:text;not null"`
	ShortCode   string `gorm:"uniqueIndex;not null"`
	ExpiresAt   *time.Time
}
