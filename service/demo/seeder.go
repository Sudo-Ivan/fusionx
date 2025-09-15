package demo

import (
	"log/slog"
	"strings"

	"github.com/0x2e/fusion/model"
	"github.com/0x2e/fusion/repo"
)

type FeedSeeder struct {
	feedRepo  *repo.Feed
	groupRepo *repo.Group
}

func NewFeedSeeder(feedRepo *repo.Feed, groupRepo *repo.Group) *FeedSeeder {
	return &FeedSeeder{
		feedRepo:  feedRepo,
		groupRepo: groupRepo,
	}
}

func (s *FeedSeeder) SeedFeeds(feedUrls string) error {
	if feedUrls == "" {
		return nil
	}

	slog.Info("Seeding demo feeds", "feeds", feedUrls)

	urls := strings.Split(feedUrls, ",")
	
	defaultGroup, err := s.groupRepo.Get(1)
	if err != nil {
		slog.Error("Failed to get default group", "error", err)
		return err
	}

	existing, err := s.feedRepo.List(nil)
	if err != nil {
		slog.Error("Failed to list existing feeds", "error", err)
		return err
	}

	existingUrls := make(map[string]bool)
	for _, feed := range existing {
		if feed.Link != nil {
			existingUrls[*feed.Link] = true
		}
	}

	var newFeeds []*model.Feed
	for i, url := range urls {
		url = strings.TrimSpace(url)
		if url == "" {
			continue
		}

		if existingUrls[url] {
			slog.Debug("Feed already exists, skipping", "url", url)
			continue
		}

		feedName := s.generateFeedName(url, i+1)
		feed := &model.Feed{
			Name:    &feedName,
			Link:    &url,
			GroupID: defaultGroup.ID,
		}

		newFeeds = append(newFeeds, feed)
		slog.Info("Prepared demo feed", "name", feedName, "url", url)
	}

	if len(newFeeds) > 0 {
		err = s.feedRepo.Create(newFeeds)
		if err != nil {
			slog.Error("Failed to create demo feeds", "error", err)
			return err
		}
		slog.Info("Successfully created demo feeds", "count", len(newFeeds))
	} else {
		slog.Info("No new demo feeds to create")
	}

	return nil
}

func (s *FeedSeeder) generateFeedName(url string, index int) string {
	url = strings.TrimPrefix(url, "https://")
	url = strings.TrimPrefix(url, "http://")
	
	parts := strings.Split(url, "/")
	if len(parts) > 0 {
		domain := parts[0]
		if domain != "" {
			return domain
		}
	}
	
	return "Demo Feed " + string(rune('A' + index - 1))
}
