package server

import (
	"context"
	"crypto/subtle"
	"fmt"
	"net/http"
	"time"

	"github.com/didip/tollbooth/v7"
	"github.com/didip/tollbooth_chi"
	"github.com/go-chi/chi/v5"
	log "github.com/go-pkgz/lgr"
	"github.com/go-pkgz/rest"

	"github.com/witjem/feedpls/internal/feed"
)

// Config the server configuration.
type Config struct {
	Listen           int
	Secret           string
	Version          string
	CacheTTL         time.Duration
	FeedsCfgFilePath string
}

// Server implement http api for gets RSS/Atom feeds.
type Server struct {
	Config
	feeds *FeedsCache
}

func NewServer(cfg Config) (*Server, error) {
	feeds, err := feed.NewFeedsFromYaml(cfg.FeedsCfgFilePath)
	if err != nil {
		return nil, fmt.Errorf("create server: %w", err)
	}

	return &Server{
		Config: cfg,
		feeds:  NewFeedsCache(feeds, cfg.CacheTTL),
	}, nil
}

// Run starts http server and closes on context cancellation.
func (s *Server) Run(ctx context.Context) error {
	s.feeds.Run(ctx)

	log.Printf("[INFO] start http server on %d", s.Listen)
	httpServer := &http.Server{
		Addr:              fmt.Sprintf(":%d", s.Listen),
		Handler:           s.router(),
		ReadHeaderTimeout: time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       time.Second,
		ErrorLog:          log.ToStdLogger(log.Default(), "WARN"),
	}

	go func() {
		<-ctx.Done()
		if httpServer != nil {
			if err := httpServer.Close(); err != nil {
				log.Printf("[ERROR] close http server, %v", err)
			}
		}

	}()

	return httpServer.ListenAndServe()
}

func (s *Server) router() http.Handler {
	router := chi.NewRouter()
	router.Use(rest.Recoverer(log.Default()))
	router.Use(rest.Throttle(100)) // limit total number of the running requests
	router.Use(rest.AppInfo("FeedPLS", "witjem", s.Version))
	router.Use(rest.Ping)
	router.Use(tollbooth_chi.LimitHandler(tollbooth.NewLimiter(10, nil)))

	router.Get("/rss/{feed}", s.getRssFeed)
	router.Get("/atom/{feed}", s.getAtomFeed)

	return router
}

// GET /rss/{feed}?secret={secret}.
func (s *Server) getRssFeed(w http.ResponseWriter, r *http.Request) {
	feedID := chi.URLParam(r, "feed")
	secret := r.URL.Query().Get("secret")

	if subtle.ConstantTimeCompare([]byte(secret), []byte(s.Secret)) != 1 {
		http.Error(w, "rejected", http.StatusForbidden)

		return
	}

	res, err := s.feeds.Get(feedID)
	if err != nil {
		http.NotFound(w, r)

		return
	}

	renderRSS(w, res)
}

// GET /atom/{feed}?secret={secret}.
func (s *Server) getAtomFeed(w http.ResponseWriter, r *http.Request) {
	feedID := chi.URLParam(r, "feed")
	secret := r.URL.Query().Get("secret")

	if subtle.ConstantTimeCompare([]byte(secret), []byte(s.Secret)) != 1 {
		http.Error(w, "rejected", http.StatusForbidden)

		return
	}

	res, err := s.feeds.Get(feedID)
	if err != nil {
		http.NotFound(w, r)

		return
	}

	renderAtom(w, res)
}
