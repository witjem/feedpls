package feed

import (
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"golang.org/x/net/html"
)

var ErrMatchNotFound = errors.New("no match found")

// Matcher defined how find data from HTML.
// Attr is optional. Init it if you want to get data from attribute.
type Matcher struct {
	Selector string `yaml:"selector"`
	Attr     string `yaml:"attr"`
}

// Find gets the first matched data.
func (m Matcher) Find(doc *goquery.Document) (string, error) {
	res, err := m.FindAll(doc)
	if err != nil {
		return "", err
	}

	return res[0], err
}

// FindAll gets all matched data.
// When no data found, return ErrMatchNotFound.
func (m Matcher) FindAll(doc *goquery.Document) ([]string, error) {
	nodes := doc.Find(m.Selector).Nodes
	res := make([]string, 0, len(nodes))

	for _, node := range nodes {
		if m.Attr != "" {
			for _, attr := range node.Attr {
				if attr.Key == m.Attr {
					res = append(res, attr.Val)

					break
				}
			}

			continue
		}

		// trys read text from HTML tag, like <title>text</title>
		if node.LastChild != nil && node.LastChild.Type == html.TextNode {
			res = append(res, node.LastChild.Data)
		}
	}

	if len(res) == 0 {
		return nil, ErrMatchNotFound
	}

	return res, nil
}

// TimeMatcher is matcher for time.
type TimeMatcher struct {
	Matcher `yaml:",inline"`

	// Layout is a standard time layout for example time.RFC3339
	Layout string `yaml:"layout"`
}

// FindTime gets matched time.
func (t TimeMatcher) FindTime(doc *goquery.Document) (time.Time, error) {
	content, err := t.Find(doc)
	if err != nil {
		return time.Time{}, err
	}

	return time.Parse(t.Layout, content)
}
