package server

import (
	"context"
	"time"

	"github.com/bluele/gcache"

	"github.com/witjem/newsfeedplease/internal/feed"
)

type FeedsCache struct {
	gc gcache.Cache
}

func NewFeedsCache(feeds Feeds, size int, ttl time.Duration) *FeedsCache {
	return &FeedsCache{
		gc: gcache.New(size).
			LoaderExpireFunc(
				func(key interface{}) (interface{}, *time.Duration, error) {
					ctx, cancel := context.WithTimeout(context.Background(), ttl)
					defer cancel()

					res, err := feeds.Get(ctx, key.(string))

					return res, &ttl, err
				},
			).
			Build(),
	}
}

func (c *FeedsCache) Get(_ context.Context, feedID string) (feed.Feed, error) {
	res, err := c.gc.Get(feedID)

	return res.(feed.Feed), err
}
