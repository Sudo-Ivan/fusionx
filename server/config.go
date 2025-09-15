package server

import (
	"context"
	"time"
)

const (
	ConfigKeyFeedRefreshInterval = "feed_refresh_interval"
	DefaultFeedRefreshInterval   = 30 * time.Minute

	ConfigKeyReadingPaneMode = "reading_pane_mode"
	DefaultReadingPaneMode   = "default" // "default", "3pane", "drawer"
)

type ConfigRepo interface {
	Get(key string) (string, error)
	Set(key, value string) error
	GetDuration(key string, defaultValue time.Duration) (time.Duration, error)
	SetDuration(key string, value time.Duration) error
}

type Config struct {
	repo     ConfigRepo
	demoMode bool
}

func NewConfig(repo ConfigRepo, demoMode bool) *Config {
	return &Config{
		repo:     repo,
		demoMode: demoMode,
	}
}

type ReqConfigUpdate struct {
	FeedRefreshIntervalMinutes int    `json:"feed_refresh_interval_minutes,omitempty" validate:"omitempty,min=1,max=10080"`
	ReadingPaneMode           string `json:"reading_pane_mode,omitempty" validate:"omitempty,oneof=default 3pane drawer"`
}

type RespConfig struct {
	FeedRefreshIntervalMinutes int    `json:"feed_refresh_interval_minutes"`
	ReadingPaneMode           string `json:"reading_pane_mode"`
	DemoMode                  bool   `json:"demo_mode"`
}

func (c *Config) Get(ctx context.Context) (*RespConfig, error) {
	interval, err := c.repo.GetDuration(ConfigKeyFeedRefreshInterval, DefaultFeedRefreshInterval)
	if err != nil {
		return nil, err
	}

	readingPaneMode, err := c.repo.Get(ConfigKeyReadingPaneMode)
	if err != nil {
		readingPaneMode = DefaultReadingPaneMode
	}

	return &RespConfig{
		FeedRefreshIntervalMinutes: int(interval.Minutes()),
		ReadingPaneMode:           readingPaneMode,
		DemoMode:                  c.demoMode,
	}, nil
}

func (c *Config) Update(ctx context.Context, req *ReqConfigUpdate) error {
	// Always update both fields when provided
	if req.FeedRefreshIntervalMinutes > 0 {
		interval := time.Duration(req.FeedRefreshIntervalMinutes) * time.Minute
		if err := c.repo.SetDuration(ConfigKeyFeedRefreshInterval, interval); err != nil {
			return err
		}
	}

	if req.ReadingPaneMode != "" {
		if err := c.repo.Set(ConfigKeyReadingPaneMode, req.ReadingPaneMode); err != nil {
			return err
		}
	}

	return nil
}

func (c *Config) GetFeedRefreshInterval() (time.Duration, error) {
	return c.repo.GetDuration(ConfigKeyFeedRefreshInterval, DefaultFeedRefreshInterval)
}
