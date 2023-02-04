package feed

import (
	"context"
	"time"
)

type Feed struct {
	ID          string
	Title       string
	Link        string
	Description string
	Items       []Item
}

type Item struct {
	ID          string
	Title       string
	Link        string
	Description string
	Published   time.Time
}

type Middleware interface {
	Process(ctx context.Context, feed Feed) (Feed, error)
}

type MiddlewareFunc func(ctx context.Context, feed Feed) (Feed, error)

func (f MiddlewareFunc) Process(ctx context.Context, feed Feed) (Feed, error) {
	return f(ctx, feed)
}

type Service struct {
	repo        *Repository
	middlewares []Middleware
}

func NewService(configs []Config, client HTTPClient, middlewares []Middleware) *Service {
	return &Service{
		repo:        NewRepository(configs, client),
		middlewares: middlewares,
	}
}

func (s *Service) Get(ctx context.Context, feedID string) (Feed, error) {
	feed, err := s.repo.Get(ctx, feedID)
	for _, middleware := range s.middlewares {
		feed, err = middleware.Process(ctx, feed)
	}

	return feed, err
}

func (s *Service) IDs() []string {
	return s.repo.IDs()
}
