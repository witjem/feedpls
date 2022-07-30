package feed

//go:generate sh -c "mockery --name=WebClient --filename=web_client.go"

import (
	"context"
	"io"
	"log"
	"net/url"
	"strings"

	"github.com/pkg/errors"

	"github.com/PuerkitoBio/goquery"
)

type Config struct {
	ID             string
	Title          string
	Description    string
	SourceURL      *url.URL
	Matchers       Matchers
	PostProcessors []PostProcessor
}

type Matchers struct {
	Page        Matcher
	Title       Matcher
	Description Matcher
	Published   TimeMatcher
}

type PostProcessor struct {
	Type    string
	Section string
	Props   map[string]string
}

type WebClient interface {
	Get(ctx context.Context, url string) (io.ReadCloser, error)
}

type Source struct {
	Config
	webClient WebClient
}

func NewSource(cfg Config, webClient WebClient) Source {
	return Source{
		Config:    cfg,
		webClient: webClient,
	}
}

// Fetch - Trying fetch Feed from source
func (s Source) Fetch(ctx context.Context) (Feed, error) {
	content, err := s.webClient.Get(ctx, s.SourceURL.String())
	if err != nil {
		return Feed{}, errors.Wrap(err, "failed to get main page content")
	}

	defer content.Close()

	doc, err := goquery.NewDocumentFromReader(content)
	if err != nil {
		return Feed{}, errors.Wrap(err, "failed to parse main page content")
	}

	matchers := s.Matchers
	links, err := matchers.Page.FindAll(doc)
	if err != nil {
		return Feed{}, errors.Wrap(err, "failed to find links from main page")
	}

	items := make([]Item, 0, len(links))
	for _, link := range links {
		item, err := s.fetchItem(ctx, link)
		if err != nil {
			log.Printf("[WARN] skip page [%s] from source [%s], err: [%s]", link, s.SourceURL, err)

			continue
		}

		items = append(items, item)
	}

	return Feed{
		ID:          s.ID,
		Title:       s.Title,
		Link:        s.SourceURL.String(),
		Description: s.Description,
		Items:       items,
	}, nil
}

// fetchItem gets Item by link
// When any item Matchers (for example title or description) not found, then return error
func (s Source) fetchItem(ctx context.Context, link string) (Item, error) {
	itemURL, err := s.toURL(link)
	if err != nil {
		return Item{}, errors.Wrap(err, "failed fetch item")
	}

	newsContent, err := s.webClient.Get(ctx, itemURL.String())
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
		Link:        itemURL.String(),
		Description: description,
		Published:   published,
	}, nil
}

func (s Source) toURL(link string) (*url.URL, error) {
	if strings.HasPrefix(link, "/") {
		return url.Parse(s.SourceURL.Scheme + "://" + s.SourceURL.Host + link)
	}

	return url.Parse(link)
}
