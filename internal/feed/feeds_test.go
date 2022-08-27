package feed_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/witjem/feedpls/internal/feed"
)

func TestCreateFeedsFromYaml(t *testing.T) {
	cfg, err := feed.ReadFeedsConfigs("./testdata/valid-feeds-config.yaml")
	assert.NoError(t, err)
	expected := feed.FeedsConfig{
		feed.SourceConfig{
			FeedID:      "first",
			Title:       "first title",
			Description: "first description",
			URL:         "https://first.example.com",
			Matchers: feed.Matchers{
				ItemURL: feed.Matcher{
					Selector: "main a",
					Attr:     "href",
				},
				Title: feed.Matcher{
					Selector: "meta[property='og:title']",
					Attr:     "content",
				},
				Description: feed.Matcher{
					Selector: "meta[property='twitter:description']",
					Attr:     "content",
				},
				Published: feed.TimeMatcher{
					Matcher: feed.Matcher{
						Selector: "meta[property='article:published_time']",
						Attr:     "content",
					},
					Layout: "2021:01:01T23:23:00",
				},
			},
		},
		feed.SourceConfig{
			FeedID:      "second",
			Title:       "second title",
			Description: "second description",
			URL:         "https://second.example.com",
			Matchers: feed.Matchers{
				ItemURL: feed.Matcher{
					Selector: "body main a",
					Attr:     "href",
				},
				Title: feed.Matcher{
					Selector: "meta[property='og:title']",
					Attr:     "content",
				},
				Description: feed.Matcher{
					Selector: "meta[property='twitter:description']",
					Attr:     "content",
				},
				Published: feed.TimeMatcher{
					Matcher: feed.Matcher{
						Selector: "meta[property='article:published_time']",
						Attr:     "content",
					},
					Layout: "2021:01:01T23:23:00",
				},
			},
		},
	}
	assert.Equal(t, expected, cfg)
}
