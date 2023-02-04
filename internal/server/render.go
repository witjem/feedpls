package server

import (
	"net/http"

	log "github.com/go-pkgz/lgr"
	gfeed "github.com/gorilla/feeds"

	"github.com/witjem/feedpls/internal/pkg/feed"
)

func renderAtom(w http.ResponseWriter, f feed.Feed) {
	w.Header().Set("Content-Type", "application/atom+xml; charset=utf-8")

	err := toGFeed(f).WriteAtom(w)
	if err != nil {
		log.Printf("[ERROR] render feed to Atom, %v", err)
		http.Error(w, "failed render feed to Atom", http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
}

func renderRSS(w http.ResponseWriter, f feed.Feed) {
	w.Header().Set("Content-Type", "application/rss+xml; charset=utf-8")

	err := toGFeed(f).WriteRss(w)
	if err != nil {
		log.Printf("[ERROR] render feed to RSS, %v", err)
		http.Error(w, "failed render feed to RSS", http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
}

func toGFeed(f feed.Feed) *gfeed.Feed {
	res := &gfeed.Feed{
		Title:       f.Title,
		Link:        &gfeed.Link{Href: f.Link},
		Description: f.Description,
		Id:          f.ID,
	}

	for _, item := range f.Items {
		res.Items = append(res.Items, &gfeed.Item{
			Id:          item.ID,
			Title:       item.Title,
			Link:        &gfeed.Link{Href: item.Link},
			Description: item.Description,
			Created:     item.Published,
		})
	}

	return res
}
