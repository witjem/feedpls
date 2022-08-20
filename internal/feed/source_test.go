package feed_test

import (
	"context"
	"io"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/witjem/newsfeedplease/internal/feed"
	"github.com/witjem/newsfeedplease/internal/feed/mocks"
)

func TestFetchFeed(t *testing.T) {
	webClient := mocks.NewWebClient(t)
	ctx := context.Background()

	webClient.On("Get", ctx, "https://blog-news.example").
		Return(file(t, "testdata/blog-news.html"), nil)
	webClient.On("Get", ctx, "https://blog-news.example/post-01").
		Return(file(t, "testdata/post-01.html"), nil)
	webClient.On("Get", ctx, "https://blog-news.example/post-02").
		Return(file(t, "testdata/post-02.html"), nil)

	source := feed.NewSource(feed.SourceConfig{
		FeedID:      "news",
		Title:       "Title",
		Description: "Description",
		URL:         "https://blog-news.example",
		Matchers: feed.Matchers{
			ItemURL: feed.Matcher{
				Selector: "main .posts a",
				Attr:     "href",
			},
			Title: feed.Matcher{
				Selector: "meta[name='twitter:title']",
				Attr:     "content",
			},
			Description: feed.Matcher{
				Selector: "meta[name='twitter:description']",
				Attr:     "content",
			},
			Published: feed.TimeMatcher{
				Matcher: feed.Matcher{
					Selector: "meta[name='article:published_time']",
					Attr:     "content",
				},
				Layout: "2006-01-02T15:04:05.000Z",
			},
		},
	}, webClient)

	actualFeed, err := source.Fetch(ctx)
	assert.NoError(t, err)

	assert.Equal(t, feed.Feed{
		ID:          "news",
		Title:       "Title",
		Link:        "https://blog-news.example",
		Description: "Description",
		Items: []feed.Item{
			{
				Title:       "Title post-01",
				Link:        "https://blog-news.example/post-01",
				Description: "Description post-01",
				Published:   time.Date(2022, 7, 22, 12, 30, 33, 0, time.UTC),
			},
			{
				Title:       "Title post-02",
				Link:        "https://blog-news.example/post-02",
				Description: "Description post-02",
				Published:   time.Date(2022, 7, 23, 12, 30, 33, 0, time.UTC),
			},
		},
	}, actualFeed)

}

func file(t *testing.T, path string) io.ReadCloser {
	t.Helper()

	f, err := os.Open(path)
	assert.NoError(t, err)

	return f
}
