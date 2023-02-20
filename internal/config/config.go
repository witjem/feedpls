package config

import (
	"github.com/witjem/feedpls/internal/pkg/funcs"
	"github.com/witjem/feedpls/internal/pkg/query"
)

type FeedsConfig []FeedConfig

type FeedConfig struct {
	FeedID      string     `yaml:"id"`
	Title       string     `yaml:"title"`
	Description string     `yaml:"description"`
	URL         string     `yaml:"url"`
	Matcher     Matcher    `yaml:"matchers"`
	Functions   []Function `yaml:"functions"`
}

type Matcher struct {
	Engine      query.Engine `yaml:"engine"`
	ItemURL     Selector     `yaml:"itemUrl"`
	Title       Selector     `yaml:"title"`
	Description Selector     `yaml:"description"`
	Published   Selector     `yaml:"published"`
}

type Selector struct {
	// required fields for GoQuery engine, options for other
	Select string `yaml:"selector"`
	Attr   string `yaml:"attr"`

	// required field for XPath engine, options for other
	Expr string `yaml:"expr"`

	// required field for time props
	Layout   string `yaml:"layout"`
	TimeZone string `yaml:"tz"`
	Locale   string `yaml:"locale"`
}

type Function struct {
	Replace *ReplaceFunc
}

type ReplaceFunc struct {
	Field funcs.Field
	From  string
	To    string
}
