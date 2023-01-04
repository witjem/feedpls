package query_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"

	"github.com/witjem/feedpls/internal/pkg/query"
)

var htmlPage = `
<!doctype html>
<html lang="uk" class="no-js" dir="ltr">
<head>	
	<meta data-react-helmet="true" name="twitter:description" content="Good news for world"/>
	<meta data-react-helmet="true" name="twitter:title" content="Good news"/>
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
		assert.Equal(t, []string{"twitter:description", "twitter:title"}, content)

		content, err = doc.Find(query.Selector{Expr: "//meta", Attr: "name"})
		assert.NoError(t, err)
		assert.Equal(t, []string{"twitter:description", "twitter:title"}, content)
	})

	t.Run("should get empty list if content by tags not found", func(t *testing.T) {
		content, err := doc.Find(query.Selector{Expr: "//a"})
		assert.NoError(t, err)
		assert.Empty(t, content)
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
		assert.Equal(t, []string{"twitter:description", "twitter:title"}, content)
	})

	t.Run("should get empty list if content by tags not found", func(t *testing.T) {
		content, err := doc.Find(query.Selector{Expr: "a"})
		assert.NoError(t, err)
		assert.Empty(t, content)
	})
}
