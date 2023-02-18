package config

import (
	"fmt"
	"time"

	"github.com/witjem/feedpls/internal/pkg/feed"
	"github.com/witjem/feedpls/internal/pkg/query"
)

func (f FeedConfig) toQueryFeedConfig() (feed.Config, error) {
	publishedTime, err := f.Matcher.Published.toQuerySelectorTime()
	if err != nil {
		return feed.Config{}, err
	}

	return feed.Config{
		FeedID:      f.FeedID,
		Title:       f.Title,
		Description: f.Description,
		URL:         f.URL,
		Matcher: feed.Matcher{
			Engine:      f.Matcher.Engine,
			ItemURL:     f.Matcher.ItemURL.toQuerySelector(),
			Title:       f.Matcher.Title.toQuerySelector(),
			Description: f.Matcher.Description.toQuerySelector(),
			Published:   publishedTime,
		},
	}, nil
}

func (s Selector) toQuerySelector() query.Selector {
	expr := s.Expr
	if s.Expr == "" {
		expr = s.Select
	}

	return query.Selector{
		Expr: expr,
		Attr: s.Attr,
	}
}

func (s Selector) toQuerySelectorTime() (query.SelectorTime, error) {
	if s.TimeZone == "" {
		s.TimeZone = "UTC"
	}

	tz, err := time.LoadLocation(s.TimeZone)
	if err != nil {
		return query.SelectorTime{}, fmt.Errorf("failed load time zone %s, %v", s.TimeZone, err)
	}

	return query.SelectorTime{
		Selector: s.toQuerySelector(),
		Layout:   s.Layout,
		TZ:       tz,
	}, nil
}
