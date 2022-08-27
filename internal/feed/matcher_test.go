package feed_test

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"

	"github.com/witjem/feedpls/internal/feed"
)

func TestMatcherFindDataFromHTMLTag(t *testing.T) {
	doc := newDoc(t, "<html><title>title text</title></html>")

	actual, err := feed.Matcher{Selector: "title"}.Find(doc)
	assert.NoError(t, err)
	assert.Equal(t, "title text", actual)
}

func TestMatcherFindDataFromHTMLTagAttribute(t *testing.T) {
	html := `<html>
				<meta name="title" content="title text">
				<meta name="description" content='title text'>
			 </html>"`
	doc := newDoc(t, html)

	actual, err := feed.Matcher{
		Selector: "meta[name='title']",
		Attr:     "content",
	}.Find(doc)
	assert.NoError(t, err)
	assert.Equal(t, "title text", actual)
}

func newDoc(t *testing.T, html string) *goquery.Document {
	t.Helper()
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	assert.NoError(t, err)

	return doc
}
