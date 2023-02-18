package query

import (
	"errors"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/antchfx/htmlquery"
	log "github.com/go-pkgz/lgr"
	"golang.org/x/net/html"
)

type Engine string

const (
	GoQuery Engine = "goquery"
	XPath   Engine = "xpath"
)

var (
	ErrEngineNotFound    = errors.New("engine not found")
	ErrCouldNotFindElem  = errors.New("could not find element in content by selector")
	ErrTimeLayoutIsEmpty = errors.New("time layout props is empty")
)

type Selector struct {
	Expr string
	Attr string
}

type SelectorTime struct {
	Selector
	Layout string
	TZ     *time.Location
}

type queryFunc func(root *html.Node, selector Selector) []*html.Node

type Document struct {
	queryFunc
	root *html.Node
}

func NewDocument(engine Engine, root *html.Node) (*Document, error) {
	var f queryFunc
	switch engine {
	case XPath:
		f = func(root *html.Node, selector Selector) []*html.Node {
			return htmlquery.Find(root, selector.Expr)
		}
	case GoQuery:
		f = func(root *html.Node, selector Selector) []*html.Node {
			doc := goquery.NewDocumentFromNode(root)

			return doc.Find(selector.Expr).Nodes
		}
	default:
		return nil, ErrEngineNotFound
	}

	return &Document{queryFunc: f, root: root}, nil
}

func (d *Document) Find(selector Selector) ([]string, error) {
	var res []string
	nodes := d.queryFunc(d.root, selector)

	for _, node := range nodes {
		if selector.Attr != "" {
			value := selectAttr(node, selector.Attr)
			if value != "" {
				res = append(res, value)
			}

			continue
		}

		if node.Type == html.TextNode {
			res = append(res, node.Data)
			continue
		}

		if node.Type == html.ElementNode {
			// trys read text from HTML tag, like <title>text</title>
			if node.LastChild != nil && node.LastChild.Type == html.TextNode {
				res = append(res, node.LastChild.Data)
			}
		}
	}

	return res, nil
}

func (d *Document) FindOne(selector Selector) (string, error) {
	res, err := d.Find(selector)
	if err != nil {
		return "", err
	}

	if len(res) < 1 {
		return "", ErrCouldNotFindElem
	}

	if len(res) > 1 {
		log.Printf("[WARN] found several elements instead of one in the by selector %+v", selector)
	}

	return res[0], nil
}

func (d *Document) FindTime(timeSelector SelectorTime) (time.Time, error) {
	if timeSelector.Layout == "" {
		return time.Time{}, ErrTimeLayoutIsEmpty
	}

	res, err := d.FindOne(timeSelector.Selector)
	if err != nil {
		return time.Time{}, err
	}

	parsedTime, err := time.ParseInLocation(timeSelector.Layout, res, timeSelector.TZ)
	if err != nil {
		return time.Time{}, err
	}

	return parsedTime, nil
}

func selectAttr(node *html.Node, attr string) string {
	return htmlquery.SelectAttr(node, attr)
}
