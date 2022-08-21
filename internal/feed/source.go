package feed

//go:generate sh -c "mockery --name=WebClient --filename=web_client.go"

import (
	"context"
	"io"
	"net/url"
	"strings"

	log "github.com/go-pkgz/lgr"

	"github.com/pkg/errors"

	"github.com/PuerkitoBio/goquery"
)

// SourceConfig describe base feed data and how gets items data.
type SourceConfig struct {
	FeedID      string   `yaml:"id"`
	Title       string   `yaml:"title"`
	Description string   `yaml:"description"`
	URL         string   `yaml:"url"`
	Matchers    Matchers `yaml:"matchers"`
}

// Matchers description of how to find data for Item.
type Matchers struct {
	ItemURL     Matcher     `yaml:"itemUrl"`
	Title       Matcher     `yaml:"title"`
	Description Matcher     `yaml:"description"`
	Published   TimeMatcher `yaml:"published"`
}

// WebClient interface for get content by URL.
type WebClient interface {
	Get(ctx context.Context, url string) (io.ReadCloser, error)
}

// Source responsible for gets Feed using by specified SourceConfig.
type Source struct {
	SourceConfig
	webClient WebClient
}

// NewSource responsible for creating Source.
func NewSource(cfg SourceConfig, webClient WebClient) *Source {
	return &Source{
		SourceConfig: cfg,
		webClient:    webClient,
	}
}

// Fetch trying fetch Feed from source.
func (s *Source) Fetch(ctx context.Context) (Feed, error) {
	content, err := s.webClient.Get(ctx, s.URL)
	if err != nil {
		return Feed{}, errors.Wrap(err, "failed to get main page content")
	}

	defer content.Close()

	doc, err := goquery.NewDocumentFromReader(content)
	if err != nil {
		return Feed{}, errors.Wrap(err, "failed to parse main page content")
	}

	matchers := s.Matchers
	links, err := matchers.ItemURL.FindAll(doc)
	if err != nil {
		return Feed{}, errors.Wrap(err, "failed to find links from main page")
	}

	items := make([]Item, 0, len(links))
	for _, link := range links {
		item, err := s.getItem(ctx, link)
		if err != nil {
			log.Printf("[WARN] skip page %s from source %s, %v", link, s.URL, err)
			continue
		}

		items = append(items, item)
	}

	return Feed{
		ID:          s.FeedID,
		Title:       s.Title,
		Link:        s.URL,
		Description: s.Description,
		Items:       items,
	}, nil
}

// getItem gets Item by link.
// When any item Matchers (for example title or description) not found, then return error.
func (s *Source) getItem(ctx context.Context, link string) (Item, error) {
	itemURL, err := s.toItemURL(link)
	if err != nil {
		return Item{}, errors.Wrap(err, "failed fetch item")
	}

	newsContent, err := s.webClient.Get(ctx, itemURL)
	if err != nil {
		return Item{}, errors.Wrap(err, "failed get item content")
	}

	defer newsContent.Close()

	newsDoc, err := goquery.NewDocumentFromReader(newsContent)
	if err != nil {
		return Item{}, errors.Wrap(err, "failed to parse item content")
	}

	matchers := s.Matchers
	title, err := matchers.Title.Find(newsDoc)
	if err != nil {
		return Item{}, errors.Wrap(err, "failed to get item title")
	}

	description, err := matchers.Description.Find(newsDoc)
	if err != nil {
		return Item{}, errors.Wrap(err, "failed to get item description")
	}

	published, err := matchers.Published.FindTime(newsDoc)
	if err != nil {
		return Item{}, errors.Wrap(err, "failed to get item published time")
	}

	return Item{
		Title:       title,
		Link:        itemURL,
		Description: description,
		Published:   published,
	}, nil
}

func (s *Source) toItemURL(link string) (string, error) {
	if strings.HasPrefix(link, "/") {
		sourceURL, err := url.Parse(s.URL)
		if err != nil {
			return "", err
		}

		return sourceURL.Scheme + "://" + sourceURL.Host + link, nil
	}

	return link, nil
}
