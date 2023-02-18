package feed_test

import (
	"context"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/witjem/feedpls/internal/pkg/feed"
	"github.com/witjem/feedpls/internal/pkg/query"
)

var feedHTML = `
<!doctype html>
<html lang="uk" class="no-js" dir="ltr">
<head>	
	<meta name="twitter:description" content="The world and galactic good news"/>
	<meta name="twitter:title" content="Good news"/>
	<title>The good news</title>
</head>
<body>
<article class="Post">
	<ul>
		<li><a href="/posts/what_is_the_best_way_to_learn_go">What is the best way to learn Go?</a></li>
		<li><a href="/posts/special_way_for_us">Special way for us</a></li>
	</ul>
</article>
</body>
`

var firstItemHTML = `
<!doctype html>
<html lang="uk" class="no-js" dir="ltr">
<head>	
	<meta name="twitter:description" content="Go is a great language. It is simple, fast, and has a large ecosystem"/>
	<meta name="twitter:title" content="What is the best way to learn Go?"/>
	<meta name="parsely-pub-date" content="2021-12-22T12:30:00Z"/>
	<title>The good news</title>
</head>
`

var secondItemHTML = `
<!doctype html>
<html lang="uk" class="no-js" dir="ltr">
<head>
	<meta name="twitter:description" content="It is a special way for us"/>
	<meta name="twitter:title" content="The special way for us"/>
	<meta name="parsely-pub-date" content="2021-12-23T14:48:00Z"/>
	<title>The good news</title>
</head>
`

func TestGetFeed(t *testing.T) {
	configs := []feed.Config{
		{
			FeedID:      "news",
			Title:       "World news",
			Description: "The world and galactic good news",
			URL:         "https://goodnews.com",
			Matcher: feed.Matcher{
				Engine: query.XPath,
				ItemURL: query.Selector{
					Expr: "//ul//a/@href",
				},
				Title: query.Selector{
					Expr: "//meta[@name='twitter:title']/@content",
				},
				Description: query.Selector{
					Expr: "//meta[@name='twitter:description']/@content",
				},
				Published: query.SelectorTime{
					Selector: query.Selector{
						Expr: "//meta[@name='parsely-pub-date']/@content",
					},
					Layout: "2006-01-02T15:04:05Z",
					TZ:     time.UTC,
				},
			},
		},
	}

	pages := map[string]string{
		"https://goodnews.com": feedHTML,
		"https://goodnews.com/posts/what_is_the_best_way_to_learn_go": firstItemHTML,
		"https://goodnews.com/posts/special_way_for_us":               secondItemHTML,
	}

	mockedHTTPClient := &HTTPClientMock{
		GetFunc: func(ctx context.Context, url string) (io.ReadCloser, error) {
			return io.NopCloser(strings.NewReader(pages[url])), nil
		},
	}

	repo := feed.NewRepository(configs, mockedHTTPClient)

	res, err := repo.Get(context.Background(), "news")
	assert.NoError(t, err)
	assert.Equal(t, feed.Feed{
		ID:          "news",
		Title:       "World news",
		Link:        "https://goodnews.com",
		Description: "The world and galactic good news",
		Items: []feed.Item{
			{
				ID:          "0de18674dc5af778dfb143a085fc3fb6", // md5 hash of the link
				Title:       "What is the best way to learn Go?",
				Link:        "https://goodnews.com/posts/what_is_the_best_way_to_learn_go",
				Description: "Go is a great language. It is simple, fast, and has a large ecosystem",
				Published:   time.Date(2021, 12, 22, 12, 30, 0, 0, time.UTC),
			},
			{
				ID:          "ff076c38e8bccb71747a2cc39fe7ecfd", // md5 hash of the link
				Title:       "The special way for us",
				Link:        "https://goodnews.com/posts/special_way_for_us",
				Description: "It is a special way for us",
				Published:   time.Date(2021, 12, 23, 14, 48, 0, 0, time.UTC),
			},
		},
	}, res)

}
