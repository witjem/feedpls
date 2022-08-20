package feed

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"github.com/witjem/newsfeedplease/internal/httpget"
)

var ErrFeedNotFound = errors.New("feed not found")

type Feed struct {
	ID          string
	Title       string
	Link        string
	Description string
	Items       []Item
}

type Item struct {
	Title       string
	Link        string
	Description string
	Published   time.Time
}

type FeedsConfig []SourceConfig

// ReadFeedsConfigs creates FeedsConfig by yaml file.
func ReadFeedsConfigs(filename string) (FeedsConfig, error) {
	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrap(err, "failed read yaml ["+filename+"]")
	}

	cfg := FeedsConfig{}
	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		return nil, errors.Wrap(err, "failed unmarshal yaml ["+filename+"]")
	}

	return cfg, nil
}

// Feeds responsible for gets Feed by ID.
type Feeds struct {
	sources map[string]*Source
}

// NewFeedsFromYaml creates Feeds by yaml file.
func NewFeedsFromYaml(filename string) (*Feeds, error) {
	configs, err := ReadFeedsConfigs(filename)
	if err != nil {
		return nil, errors.Wrap(err, "failed create Feeds from yaml")
	}

	return NewFeeds(configs), nil
}

// NewFeeds creates Feeds by FeedsConfig.
func NewFeeds(configs FeedsConfig) *Feeds {
	sources := make(map[string]*Source)
	for _, cfg := range configs {
		sources[cfg.FeedID] = NewSource(cfg, httpget.New())
	}

	return &Feeds{sources: sources}
}

// Get gets Feed by feedID.
func (s *Feeds) Get(ctx context.Context, feedID string) (Feed, error) {
	source, ok := s.sources[feedID]
	if !ok {
		return Feed{}, ErrFeedNotFound
	}

	res, err := source.Fetch(ctx)
	if err != nil {
		return Feed{}, errors.Wrap(err, fmt.Sprintf("failed fetch feed %s", feedID))
	}

	return res, nil
}

// Size the number of registrants feed sources.
func (s *Feeds) Size() int {
	return len(s.sources)
}
