package postprocessor

import (
	"context"
	"github.com/witjem/feedpls/internal/feed"
	"github.com/witjem/feedpls/internal/server"
)

type PostProcessor struct {
	server.FeedsRepo
	cfg feed.FeedsConfig
}

func NewPostProcessor(feeds server.FeedsRepo, cfg feed.FeedsConfig) *PostProcessor {
	return &PostProcessor{FeedsRepo: feeds, cfg: cfg}
}

func (p *PostProcessor) Get(ctx context.Context, feedID string) (feed.Feed, error) {
	res, err := p.FeedsRepo.Get(ctx, feedID)
	if err != nil {
		return res, err
	}
	p.TryToReplace(&res)
	return feed.Feed{}, err
}

// TryToReplace - make replace if replace postProcessor is declared in config
func (p *PostProcessor) TryToReplace(feed *feed.Feed) {
	for _, config := range p.cfg {
		if config.FeedID != feed.ID {
			continue
		}
		if len(config.PostProcessors) == 0 {
			continue
		}
		p.ApplyReplace(&config.PostProcessors, feed)
	}
}

func (p *PostProcessor) ApplyReplace(postProcessors *[]feed.PostProcessor, feed *feed.Feed) {
	for _, postProcessor := range *postProcessors {
		if postProcessor.Replace.Field == "" {
			continue
		}
		switch postProcessor.Replace.Field {
		case "title":
			feed.Title = postProcessor.Replace.To
		case "description":
			feed.Description = postProcessor.Replace.To
		}
	}
}
