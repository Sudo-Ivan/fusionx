package pull

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/0x2e/fusion/model"
	"github.com/0x2e/fusion/pkg/ptr"
	"github.com/0x2e/fusion/repo"
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
}

// TODO: cache favicon

func NewPuller(feedRepo FeedRepo, itemRepo ItemRepo, configRepo ConfigRepo) *Puller {
	return &Puller{
		feedRepo:   feedRepo,
		itemRepo:   itemRepo,
		configRepo: configRepo,
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
