package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/witjem/feedpls/internal/pkg/feed"
	"github.com/witjem/feedpls/internal/pkg/funcs"
	"github.com/witjem/feedpls/internal/pkg/httpget"
)

// NewFeedServiceFromYML creates a new feed service from a YAML config file.
func NewFeedServiceFromYML(cfgPath string) (*feed.Service, error) {
	feedConfigs, err := ReadFeedsConfigFromYaml(cfgPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	configs := make([]feed.Config, len(feedConfigs))
	var middlewares []feed.Middleware
	for i, f := range feedConfigs {
		configs[i] = f.toQueryFeedConfig()
		middlewares = append(middlewares, toMiddlewares(f.Functions)...)
	}

	return feed.NewService(configs, httpget.New(), middlewares), nil
}

// ReadFeedsConfigFromYaml reads a YAML config file and returns a slice of FeedConfig.
func ReadFeedsConfigFromYaml(cfgPath string) (FeedsConfig, error) {
	bytes, err := os.ReadFile(cfgPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	feedConfigs := FeedsConfig{}
	err = yaml.Unmarshal(bytes, &feedConfigs)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file: %w", err)
	}

	return feedConfigs, nil
}

func toMiddlewares(functions []Function) []feed.Middleware {
	var res []feed.Middleware
	for _, function := range functions {
		if function.Replace != nil {
			replaceFunc := function.Replace
			res = append(res, funcs.Replace(replaceFunc.Field, replaceFunc.From, replaceFunc.To))
		}
	}

	return res
}
