package query_test

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"

	"github.com/goodsign/monday"

	"github.com/witjem/feedpls/internal/pkg/query"
)

var htmlPage = `
<!doctype html>
<html lang="uk" class="no-js" dir="ltr">
<head>	
	<meta data-react-helmet="true" name="twitter:description" content="Good news for world"/>
	<meta data-react-helmet="true" name="twitter:title" content="Good news"/>
	<meta name="content-date" content="January 2, 2022, 3:04 pm"/>
	<date name="date-in-uk-UA-locale" content="Січень 2, 2022, 15:04"/>
	<title>Good news</title>
</head>
`

func TestXPath(t *testing.T) {

	root, err := html.Parse(strings.NewReader(htmlPage))
	assert.NoError(t, err)

	doc, err := query.NewDocument(query.XPath, root)
	assert.NoError(t, err)

	t.Run("should get content from tag", func(t *testing.T) {
		content, err := doc.FindOne(query.Selector{Expr: "//title"})
		assert.NoError(t, err)
		assert.Equal(t, "Good news", content)
	})

	t.Run("should get content from tag attribute", func(t *testing.T) {
		content, err := doc.FindOne(query.Selector{Expr: "//meta[@name='twitter:description']/@content"})
		assert.NoError(t, err)
		assert.Equal(t, "Good news for world", content)

		content, err = doc.FindOne(query.Selector{Expr: "//meta[@name='twitter:description']", Attr: "content"})
		assert.NoError(t, err)
		assert.Equal(t, "Good news for world", content)
	})

	t.Run("should get list contents by tags", func(t *testing.T) {
		content, err := doc.Find(query.Selector{Expr: "//meta/@name"})
		assert.NoError(t, err)
		assert.Equal(t, []string{"twitter:description", "twitter:title", "content-date"}, content)

		content, err = doc.Find(query.Selector{Expr: "//meta", Attr: "name"})
		assert.NoError(t, err)
		assert.Equal(t, []string{"twitter:description", "twitter:title", "content-date"}, content)
	})

	t.Run("should get empty list if content by tags not found", func(t *testing.T) {
		content, err := doc.Find(query.Selector{Expr: "//a"})
		assert.NoError(t, err)
		assert.Empty(t, content)
	})

	t.Run("should parse time with timezone", func(t *testing.T) {
		timeLocation := time.FixedZone("UTC+3", 3*60*60)

		actual, err := doc.FindTime(
			query.SelectorTime{
				Selector: query.Selector{Expr: "//meta[@name='content-date']/@content"},
				Layout:   "January 2, 2006, 3:04 pm",
				TZ:       timeLocation,
				Locale:   monday.LocaleEnGB,
			})
		assert.NoError(t, err)
		assert.Equal(t, time.Date(2022, 1, 2, 15, 4, 0, 0, timeLocation), actual)
	})

	t.Run("should parse time with timezone and uk_UA locale", func(t *testing.T) {
		zone := time.FixedZone("UTC+3", 3*60*60)
		actual, err := doc.FindTime(
			query.SelectorTime{
				Selector: query.Selector{Expr: "//date[@name='date-in-uk-UA-locale']/@content"},
				Layout:   "January 2, 2006, 15:04",
				TZ:       zone,
				Locale:   monday.LocaleUkUA,
			})
		assert.NoError(t, err)
		assert.Equal(t, time.Date(2022, 1, 2, 15, 4, 0, 0, zone), actual)
	})
}

func TestGoQuery(t *testing.T) {

	root, err := html.Parse(strings.NewReader(htmlPage))
	assert.NoError(t, err)

	doc, err := query.NewDocument(query.GoQuery, root)
	assert.NoError(t, err)

	t.Run("should get content from tag", func(t *testing.T) {
		content, err := doc.FindOne(query.Selector{Expr: "title"})
		assert.NoError(t, err)
		assert.Equal(t, "Good news", content)
	})

	t.Run("should get content from tag attribute", func(t *testing.T) {
		content, err := doc.FindOne(query.Selector{Expr: "meta[name='twitter:description']", Attr: "content"})
		assert.NoError(t, err)
		assert.Equal(t, "Good news for world", content)
	})

	t.Run("should get list contents by tags", func(t *testing.T) {
		content, err := doc.Find(query.Selector{Expr: "meta", Attr: "name"})
		assert.NoError(t, err)
		assert.Equal(t, []string{"twitter:description", "twitter:title", "content-date"}, content)
	})

	t.Run("should get empty list if content by tags not found", func(t *testing.T) {
		content, err := doc.Find(query.Selector{Expr: "a"})
		assert.NoError(t, err)
		assert.Empty(t, content)
	})

	t.Run("should parse time with timezone", func(t *testing.T) {
		timeLocation := time.FixedZone("UTC+3", 3*60*60)

		actual, err := doc.FindTime(
			query.SelectorTime{
				Selector: query.Selector{Expr: "meta[name='content-date']", Attr: "content"},
				Layout:   "January 2, 2006, 3:04 pm",
				TZ:       timeLocation,
				Locale:   monday.LocaleEnGB,
			})
		assert.NoError(t, err)
		assert.Equal(t, time.Date(2022, 1, 2, 15, 4, 0, 0, timeLocation), actual)
	})
}
