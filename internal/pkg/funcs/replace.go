package funcs

import (
	"context"
	"strings"

	"github.com/witjem/feedpls/internal/pkg/feed"
)

type Field string

const (
	FieldTitle       = "title"
	FieldDescription = "description"
)

func Replace(field Field, from, to string) feed.MiddlewareFunc {
	return func(ctx context.Context, f feed.Feed) (feed.Feed, error) {
		for i := range f.Items {
			if field == FieldTitle {
				f.Items[i].Title = strings.Replace(f.Items[i].Title, from, to, 1)
			}

			if field == FieldDescription {
				f.Items[i].Description = strings.Replace(f.Items[i].Description, from, to, 1)
			}
		}

		return f, nil
	}
}
