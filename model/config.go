package model

import (
	"time"
)

type Config struct {
	ID        uint      `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Key   string `gorm:"key;uniqueIndex;not null"`
	Value string `gorm:"value;not null"`
}
