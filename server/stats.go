package server

import (
	"context"
	"os"
	"time"
)

type StatsRepo interface {
	GetTotalFeeds() (int, error)
	GetTotalItems() (int, error)
	GetTotalUnreadItems() (int, error)
	GetTotalGroups() (int, error)
	GetLastFeedUpdate() (*time.Time, error)
	GetFailedFeeds() (int, error)
}

type Stats struct {
	repo     StatsRepo
	dbPath   string
}

func NewStats(repo StatsRepo, dbPath string) *Stats {
	return &Stats{
		repo:   repo,
		dbPath: dbPath,
	}
}

func (s Stats) Get(ctx context.Context) (*RespStats, error) {
	totalFeeds, err := s.repo.GetTotalFeeds()
	if err != nil {
		return nil, err
	}

	totalItems, err := s.repo.GetTotalItems()
	if err != nil {
		return nil, err
	}

	totalUnreadItems, err := s.repo.GetTotalUnreadItems()
	if err != nil {
		return nil, err
	}

	totalGroups, err := s.repo.GetTotalGroups()
	if err != nil {
		return nil, err
	}

	lastFeedUpdate, err := s.repo.GetLastFeedUpdate()
	if err != nil {
		return nil, err
	}

	failedFeeds, err := s.repo.GetFailedFeeds()
	if err != nil {
		return nil, err
	}

	var dbSize int64
	if s.dbPath != "" {
		if stat, err := os.Stat(s.dbPath); err == nil {
			dbSize = stat.Size()
		}
	}

	return &RespStats{
		TotalFeeds:       totalFeeds,
		TotalItems:       totalItems,
		TotalUnreadItems: totalUnreadItems,
		TotalGroups:      totalGroups,
		DatabaseSize:     dbSize,
		LastFeedUpdate:   lastFeedUpdate,
		FailedFeeds:      failedFeeds,
	}, nil
}
