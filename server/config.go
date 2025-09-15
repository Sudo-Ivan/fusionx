package server

import (
	"context"
	"time"
)

const (
	ConfigKeyFeedRefreshInterval = "feed_refresh_interval"
	DefaultFeedRefreshInterval   = 30 * time.Minute
)

type ConfigRepo interface {
	Get(key string) (string, error)
	Set(key, value string) error
	GetDuration(key string, defaultValue time.Duration) (time.Duration, error)
	SetDuration(key string, value time.Duration) error
}

type Config struct {
	repo ConfigRepo
}

func NewConfig(repo ConfigRepo) *Config {
	return &Config{
		repo: repo,
	}
}

type ReqConfigUpdate struct {
	FeedRefreshIntervalMinutes int `json:"feed_refresh_interval_minutes" validate:"min=1,max=10080"`
}

type RespConfig struct {
	FeedRefreshIntervalMinutes int `json:"feed_refresh_interval_minutes"`
}

func (c *Config) Get(ctx context.Context) (*RespConfig, error) {
	interval, err := c.repo.GetDuration(ConfigKeyFeedRefreshInterval, DefaultFeedRefreshInterval)
	if err != nil {
		return nil, err
	}

	return &RespConfig{
		FeedRefreshIntervalMinutes: int(interval.Minutes()),
	}, nil
}

func (c *Config) Update(ctx context.Context, req *ReqConfigUpdate) error {
	interval := time.Duration(req.FeedRefreshIntervalMinutes) * time.Minute
	return c.repo.SetDuration(ConfigKeyFeedRefreshInterval, interval)
}

func (c *Config) GetFeedRefreshInterval() (time.Duration, error) {
	return c.repo.GetDuration(ConfigKeyFeedRefreshInterval, DefaultFeedRefreshInterval)
}
