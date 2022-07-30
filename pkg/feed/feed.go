package feed

import "time"

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
