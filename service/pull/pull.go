package pull

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/Sudo-Ivan/fusionx/model"
	"github.com/Sudo-Ivan/fusionx/pkg/ptr"
	"github.com/Sudo-Ivan/fusionx/repo"
	"github.com/Sudo-Ivan/fusionx/service/favicon"
)

var (
	interval = 30 * time.Minute
)

type FeedRepo interface {
	List(filter *repo.FeedListFilter) ([]*model.Feed, error)
	Get(id uint) (*model.Feed, error)
	Update(id uint, feed *model.Feed) error
}

type ItemRepo interface {
	Insert(items []*model.Item) error
}

type ConfigRepo interface {
	GetFeedRefreshInterval() (time.Duration, error)
}

type Puller struct {
	feedRepo   FeedRepo
	itemRepo   ItemRepo
	configRepo ConfigRepo
	faviconSvc *favicon.Service
}

// TODO: cache favicon

func NewPuller(feedRepo FeedRepo, itemRepo ItemRepo, configRepo ConfigRepo) *Puller {
	return &Puller{
		feedRepo:   feedRepo,
		itemRepo:   itemRepo,
		configRepo: configRepo,
		faviconSvc: favicon.NewService("./cache/favicons"),
	}
}

func (p *Puller) Run() {
	// Get initial interval
	currentInterval := p.getCurrentInterval()
	ticker := time.NewTicker(currentInterval)
	defer ticker.Stop()

	for {
		// #nosec G104 - PullAll errors are logged internally, service should continue running
		p.PullAll(context.Background(), false)
		
		// Also try to fix missing favicons
		p.FixMissingFavicons(context.Background())

		<-ticker.C

		// Check if interval has changed and update ticker if needed
		newInterval := p.getCurrentInterval()
		if newInterval != currentInterval {
			currentInterval = newInterval
			ticker.Stop()
			ticker = time.NewTicker(currentInterval)
		}
	}
}

func (p *Puller) getCurrentInterval() time.Duration {
	if p.configRepo == nil {
		return interval
	}
	
	configInterval, err := p.configRepo.GetFeedRefreshInterval()
	if err != nil {
		slog.Warn("failed to get feed refresh interval from config, using default", "error", err)
		return interval
	}
	return configInterval
}

func (p *Puller) FixMissingFavicons(ctx context.Context) {
	feeds, err := p.feedRepo.List(&repo.FeedListFilter{})
	if err != nil {
		slog.Warn("failed to get feeds for favicon fixing", "error", err)
		return
	}
	
	for _, feed := range feeds {
		if feed.Link != nil && (feed.FaviconPath == nil || *feed.FaviconPath == "") {
			// This feed doesn't have a cached favicon, try to fetch it
			if faviconPath, err := p.faviconSvc.GetFaviconPath(*feed.Link); err == nil {
				// Update the feed with the favicon path
				// #nosec G104 - favicon update is non-critical, error can be ignored
				_ = p.feedRepo.Update(feed.ID, &model.Feed{FaviconPath: &faviconPath})
				slog.Debug("fixed missing favicon", "feed_id", feed.ID, "favicon_path", faviconPath)
			}
		}
	}
}

func (p *Puller) PullAll(ctx context.Context, force bool) error {
	currentInterval := p.getCurrentInterval()
	ctx, cancel := context.WithTimeout(ctx, currentInterval/2)
	defer cancel()

	feeds, err := p.feedRepo.List(nil)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			err = nil
		}
		return err
	}
	if len(feeds) == 0 {
		return nil
	}

	routinePool := make(chan struct{}, 10)
	defer close(routinePool)
	wg := sync.WaitGroup{}
	for _, f := range feeds {
		routinePool <- struct{}{}
		wg.Add(1)
		go func(f *model.Feed) {
			defer func() {
				wg.Done()
				<-routinePool
			}()

			if err := p.do(ctx, f, force); err != nil {
				slog.Error("failed to pull feed", "error", err, "feed_id", f.ID, "feed_link", ptr.From(f.Link))
			}
		}(f)
	}
	wg.Wait()
	return nil
}

func (p *Puller) PullOne(ctx context.Context, id uint) error {
	f, err := p.feedRepo.Get(id)
	if err != nil {
		return err
	}

	return p.do(ctx, f, true)
}
