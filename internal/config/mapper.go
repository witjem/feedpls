package config

import (
	"github.com/witjem/feedpls/internal/pkg/feed"
	"github.com/witjem/feedpls/internal/pkg/query"
)

func (f FeedConfig) toQueryFeedConfig() feed.Config {
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
			Published:   f.Matcher.Published.toQuerySelectorTime(),
		},
	}
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

func (s Selector) toQuerySelectorTime() query.SelectorTime {
	return query.SelectorTime{
		Selector: s.toQuerySelector(),
		Layout:   s.Layout,
	}
}
