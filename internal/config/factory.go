package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"

	feed2 "github.com/witjem/feedpls/internal/pkg/feed"
	"github.com/witjem/feedpls/internal/pkg/funcs"
	"github.com/witjem/feedpls/internal/pkg/httpget"
)

func NewFeedServiceFromYML(cfgPath string) (*feed2.Service, error) {
	feedConfigs, err := ReadFeedsConfigFromYaml(cfgPath)
	if err != nil {
		return nil, err
	}

	configs := make([]feed2.Config, len(feedConfigs))
	var middlewares []feed2.Middleware
	for i, f := range feedConfigs {
		feedCfg := f.toQueryFeedConfig()

		err = feedCfg.Validate()
		if err != nil {
			return nil, fmt.Errorf("does not valid config for feed %s: %v", f.FeedID, err)
		}

		configs[i] = feedCfg
		middlewares = append(middlewares, toMiddlewares(f.Functions)...)
	}

	return feed2.NewService(configs, httpget.New(), middlewares), nil
}

func ReadFeedsConfigFromYaml(cfgPath string) (FeedsConfig, error) {
	bytes, err := os.ReadFile(cfgPath)
	if err != nil {
		return nil, err
	}

	feedConfigs := FeedsConfig{}
	err = yaml.Unmarshal(bytes, &feedConfigs)
	if err != nil {
		return nil, err
	}

	return feedConfigs, nil
}

func toMiddlewares(functions []Function) []feed2.Middleware {
	var res []feed2.Middleware
	for _, function := range functions {
		if function.Replace != nil {
			replaceFunc := function.Replace
			res = append(res, funcs.Replace(replaceFunc.Field, replaceFunc.From, replaceFunc.To))
		}
	}

	return res
}
