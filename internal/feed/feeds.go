package feed

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	yaml "gopkg.in/yaml.v2"

	"github.com/witjem/feedpls/internal/httpget"
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
		return nil, fmt.Errorf("find config file: %w", err)
	}

	cfg := FeedsConfig{}
	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		return nil, fmt.Errorf("parse config file: %w", err)
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
		return nil, fmt.Errorf("new feeds from yaml: %v", err)
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

	res, err := source.Get(ctx)
	if err != nil {
		return Feed{}, fmt.Errorf("get feed by id %s: %v", feedID, err)
	}

	return res, nil
}

// IDs all registered feeds ids.
func (s *Feeds) IDs() []string {
	ids := make([]string, len(s.sources))
	i := 0
	for k := range s.sources {
		ids[i] = k
		i++
	}

	return ids
}
