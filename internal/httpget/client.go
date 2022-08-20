package httpget

import (
	"context"
	"io"
	"net/http"
)

const userAgent = "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/45.0.2454.101 Safari/537.36"

// Client the simple http client.
type Client struct {
	client *http.Client
}

func New() *Client {
	return &Client{client: &http.Client{}}
}

// Get gets content by URL.
func (q Client) Get(ctx context.Context, url string) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, http.NoBody)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", userAgent)
	resp, err := q.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}
