package server

import (
	"context"
	"time"

	"github.com/bluele/gcache"
	"github.com/go-pkgz/lgr"

	"github.com/witjem/newsfeedplease/internal/feed"
)

// FeedsRepo interface for gets the latest up-to-date feeds.
type FeedsRepo interface {
	Get(ctx context.Context, feedID string) (feed.Feed, error)
	IDs() []string
}

// FeedsCache - cache for the latest up-to-date feeds.
type FeedsCache struct {
	gc     gcache.Cache
	ticker *time.Ticker
	feeds  FeedsRepo
}

// NewFeedsCache creates FeedsCache.
// ttl - duration how often will be refresh cache for gets up-to-date feeds.
func NewFeedsCache(feeds FeedsRepo, ttl time.Duration) *FeedsCache {
	return &FeedsCache{
		gc:     gcache.New(len(feeds.IDs())).Build(),
		ticker: time.NewTicker(ttl),
		feeds:  feeds,
	}
}

// Run loads actual feeds to cache and runs background tasks for updates its.
// When ctx.Done -> all background tasks stopping too.
func (f *FeedsCache) Run(ctx context.Context) {
	lgr.Printf("[INFO] start caching feeds")

	loadFeedToCache := func(feedID string) {
		lgr.Printf("[INFO] loading feed %s to cache", feedID)

		newFeed, err := f.feeds.Get(context.TODO(), feedID)
		if err != nil {
			lgr.Printf("[ERROR] load feed %s to cache: %v", feedID, err)
		}

		err = f.gc.Set(feedID, newFeed)
		if err != nil {
			lgr.Printf("[ERROR] load feed %s to cache: %v", feedID, err)
		}
	}

	refreshCache := func(ctx context.Context, feedID string) {
		for {
			select {
			case <-f.ticker.C:
				loadFeedToCache(feedID)
			case <-ctx.Done():
				f.ticker.Stop()
				lgr.Printf("[INFO] stopped refreshing cache for the feed %s", feedID)

				return
			}
		}
	}

	// load all feeds to cache
	for _, feedID := range f.feeds.IDs() {
		loadFeedToCache(feedID)
	}
	lgr.Printf("[INFO] all feeds loaded to cache")

	// run background tasks for refreshing feeds on the cache
	for _, feedID := range f.feeds.IDs() {
		go refreshCache(ctx, feedID)
	}
}

func (f *FeedsCache) Get(feedID string) (feed.Feed, error) {
	res, err := f.gc.Get(feedID)

	return res.(feed.Feed), err
}
