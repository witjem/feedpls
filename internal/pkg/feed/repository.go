package feed

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"

	log "github.com/go-pkgz/lgr"
	"golang.org/x/net/html"

	"github.com/witjem/feedpls/internal/pkg/query"
)

var ErrIDNotFound = errors.New("feed by id not found")

type Config struct {
	FeedID      string
	Title       string
	Description string
	URL         string
	Matcher     Matcher
}

type Matcher struct {
	Engine      query.Engine
	ItemURL     query.Selector
	Title       query.Selector
	Description query.Selector
	Published   query.SelectorTime
}

func (m Matcher) FindItemURLs(doc *query.Document) ([]string, error) {
	return doc.Find(m.ItemURL)
}

func (m Matcher) FindItem(doc *query.Document) (Item, error) {
	var err error
	res := Item{}

	res.Title, err = doc.FindOne(m.Title)
	if err != nil {
		return Item{}, fmt.Errorf("failed to find item title: %w", err)
	}

	res.Description, err = doc.FindOne(m.Description)
	if err != nil {
		return Item{}, fmt.Errorf("failed to find item description: %w", err)
	}

	res.Published, err = doc.FindTime(m.Published)
	if err != nil {
		return Item{}, fmt.Errorf("find item published time: %w", err)
	}

	return res, nil
}

func (c Config) Validate() error {
	// todo
	return nil
}

// HTTPClient interface for get content by URL.
type HTTPClient interface {
	Get(ctx context.Context, url string) (io.ReadCloser, error)
}

type Repository struct {
	cfg        map[string]Config
	httpClient HTTPClient
}

func NewRepository(configs []Config, client HTTPClient) *Repository {
	mCfg := make(map[string]Config)
	for _, cfg := range configs {
		mCfg[cfg.FeedID] = cfg
	}

	return &Repository{
		cfg:        mCfg,
		httpClient: client,
	}
}

func (r *Repository) Get(ctx context.Context, feedID string) (Feed, error) {
	cfg, ok := r.cfg[feedID]
	if !ok {
		return Feed{}, ErrIDNotFound
	}

	content, err := r.httpClient.Get(ctx, cfg.URL)
	if err != nil {
		return Feed{}, fmt.Errorf("get feed content: %w", err)
	}
	defer content.Close()

	nodes, err := html.Parse(content)
	if err != nil {
		return Feed{}, fmt.Errorf("get feed content: %w", err)
	}

	doc, err := query.NewDocument(cfg.Matcher.Engine, nodes)
	if err != nil {
		return Feed{}, fmt.Errorf("get feed content: %w", err)
	}

	itemURLs, err := cfg.Matcher.FindItemURLs(doc)
	if err != nil {
		return Feed{}, fmt.Errorf("find items urls: %w", err)
	}

	items := make([]Item, 0, len(itemURLs))
	for _, link := range itemURLs {

		itemURL, err := prepareURL(cfg.URL, link)
		if err != nil {
			log.Printf("[WARN] skip item %s from feed %s: validate item url: %v", itemURL, cfg.FeedID, err)
			continue
		}

		item, err := r.getItem(ctx, cfg.Matcher, itemURL)
		if err != nil {
			log.Printf("[WARN] skip item %s from feed %s: %v", itemURL, cfg.FeedID, err)
			continue
		}

		items = append(items, item)
	}

	return Feed{
		ID:          cfg.FeedID,
		Title:       cfg.Title,
		Link:        cfg.URL,
		Description: cfg.Description,
		Items:       items,
	}, nil
}

func (r *Repository) getItem(ctx context.Context, matcher Matcher, itemURL string) (Item, error) {
	content, err := r.httpClient.Get(ctx, itemURL)
	if err != nil {
		return Item{}, fmt.Errorf("get item content: %w", err)
	}
	defer content.Close()

	nodes, err := html.Parse(content)
	if err != nil {
		return Item{}, fmt.Errorf("get item content: %w", err)
	}

	doc, err := query.NewDocument(matcher.Engine, nodes)
	if err != nil {
		return Item{}, fmt.Errorf("get item content: %w", err)
	}

	return matcher.FindItem(doc)
}

func (r *Repository) IDs() []string {
	res := make([]string, len(r.cfg))
	for id := range r.cfg {
		res = append(res, id)
	}

	return res
}

func prepareURL(baseURL, itemURL string) (string, error) {
	if strings.HasPrefix(itemURL, "/") {
		sourceURL, err := url.Parse(baseURL)
		if err != nil {
			return "", err
		}

		return sourceURL.Scheme + "://" + sourceURL.Host + itemURL, nil
	}

	return itemURL, nil
}
