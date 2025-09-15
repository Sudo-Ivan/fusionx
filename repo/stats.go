package repo

import (
	"errors"
	"time"

	"github.com/Sudo-Ivan/fusionx/model"
	"gorm.io/gorm"
)

type Stats struct {
	db *gorm.DB
}

func NewStats(db *gorm.DB) *Stats {
	return &Stats{db: db}
}

func (s Stats) GetTotalFeeds() (int, error) {
	var count int64
	err := s.db.Model(&model.Feed{}).Count(&count).Error
	return int(count), err
}

func (s Stats) GetTotalItems() (int, error) {
	var count int64
	err := s.db.Model(&model.Item{}).Count(&count).Error
	return int(count), err
}

func (s Stats) GetTotalUnreadItems() (int, error) {
	var count int64
	err := s.db.Model(&model.Item{}).Where("unread = ?", true).Count(&count).Error
	return int(count), err
}

func (s Stats) GetTotalGroups() (int, error) {
	var count int64
	err := s.db.Model(&model.Group{}).Count(&count).Error
	return int(count), err
}

func (s Stats) GetLastFeedUpdate() (*time.Time, error) {
	var feed model.Feed
	err := s.db.Model(&model.Feed{}).Order("updated_at DESC").First(&feed).Error
	if err != nil {
		// For stats, no feeds is not an error, just return nil
		if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, ErrNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &feed.UpdatedAt, nil
}

func (s Stats) GetFailedFeeds() (int, error) {
	var count int64
	err := s.db.Model(&model.Feed{}).Where("failure != '' AND failure IS NOT NULL").Count(&count).Error
	return int(count), err
}
