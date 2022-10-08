package feed

//go:generate sh -c "mockery --name=WebClient --filename=web_client.go"

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	log "github.com/go-pkgz/lgr"
)

// SourceConfig describe base feed data and how gets items data.
type SourceConfig struct {
	FeedID         string          `yaml:"id"`
	Title          string          `yaml:"title"`
	Description    string          `yaml:"description"`
	URL            string          `yaml:"url"`
	Matchers       Matchers        `yaml:"matchers"`
	PostProcessors []PostProcessor `yaml:"postprocessors"`
}

// Matchers description of how to find data for Item.
type Matchers struct {
	ItemURL     Matcher     `yaml:"itemUrl"`
	Title       Matcher     `yaml:"title"`
	Description Matcher     `yaml:"description"`
	Published   TimeMatcher `yaml:"published"`
}

type PostProcessor struct {
	Replace Replace `yaml:"replace"`
}
type Replace struct {
	Field string `yaml:"field"` // can be: title, description
	From  string `yaml:"from"`  // some_regex
	To    string `yaml:"to"`
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

// Get trying get Feed from source.
func (s *Source) Get(ctx context.Context) (Feed, error) {
	content, err := s.webClient.Get(ctx, s.URL)
	if err != nil {
		return Feed{}, fmt.Errorf("get feed content: %w", err)
	}

	defer content.Close()

	doc, err := goquery.NewDocumentFromReader(content)
	if err != nil {
		return Feed{}, fmt.Errorf("parse content: %w", err)
	}

	matchers := s.Matchers
	links, err := matchers.ItemURL.FindAll(doc)
	if err != nil {
		return Feed{}, fmt.Errorf("find items urls: %w", err)
	}

	items := make([]Item, 0, len(links))
	for _, link := range links {
		item, err := s.getItem(ctx, link)
		if err != nil {
			log.Printf("[WARN] skip item %s from feed %s: %v", link, s.FeedID, err)
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
		return Item{}, fmt.Errorf("validate item url: %w", err)
	}

	newsContent, err := s.webClient.Get(ctx, itemURL)
	if err != nil {
		return Item{}, fmt.Errorf("get item content: %w", err)
	}

	defer newsContent.Close()

	newsDoc, err := goquery.NewDocumentFromReader(newsContent)
	if err != nil {
		return Item{}, fmt.Errorf("parse item content: %w", err)
	}

	matchers := s.Matchers
	title, err := matchers.Title.Find(newsDoc)
	if err != nil {
		return Item{}, fmt.Errorf("find item title: %w", err)
	}

	description, err := matchers.Description.Find(newsDoc)
	if err != nil {
		return Item{}, fmt.Errorf("find item description: %w", err)
	}

	published, err := matchers.Published.FindTime(newsDoc)
	if err != nil {
		return Item{}, fmt.Errorf("find item published time: %w", err)
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
