package postprocessor

import (
	"context"
	log "github.com/go-pkgz/lgr"
	"github.com/witjem/feedpls/internal/feed"
	"regexp"
)

type FeedsRepo interface {
	Get(ctx context.Context, feedID string) (feed.Feed, error)
	IDs() []string
}

type PostProcessor struct {
	FeedsRepo
	cfg feed.FeedsConfig
}

func NewPostProcessor(feeds FeedsRepo, cfg feed.FeedsConfig) *PostProcessor {
	return &PostProcessor{FeedsRepo: feeds, cfg: cfg}
}

func (p *PostProcessor) Get(ctx context.Context, feedID string) (feed.Feed, error) {
	res, err := p.FeedsRepo.Get(ctx, feedID)
	if err != nil {
		return res, err
	}
	p.TryToReplace(&res)
	return res, err
}

// TryToReplace - make replace if replace postProcessor is declared in config.
func (p *PostProcessor) TryToReplace(feedData *feed.Feed) {
	for _, config := range p.cfg {
		if config.FeedID != feedData.ID {
			continue
		}
		if len(config.PostProcessors) == 0 {
			continue
		}
		p.ApplyReplace(&config.PostProcessors, feedData)
	}
}

func (p *PostProcessor) ApplyReplace(postProcessors *[]feed.PostProcessor, feedData *feed.Feed) {
	for _, postProcessor := range *postProcessors {
		if postProcessor.Replace.Field == "" {
			continue
		}
		re, err := regexp.Compile(postProcessor.Replace.From)
		if err != nil {
			log.Printf("[ERROR] incorrect regexp expression, err: %v", err)
			return
		}
		switch postProcessor.Replace.Field {
		case "title":
			feedData.Title = re.ReplaceAllString(feedData.Title, postProcessor.Replace.To)
		case "description":
			feedData.Description = re.ReplaceAllString(feedData.Description, postProcessor.Replace.To)
		}
	}
}
