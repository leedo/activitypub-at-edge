package render

import (
	"io"

	"github.com/leedo/activitypub-at-edge/activitypub"
)

func Pagination(w io.Writer, c *activitypub.Collection) {
	w.Write([]byte(`<A HREF="/` + c.Prev() + `">Prev</A> | `))
	w.Write([]byte(`<A HREF="/` + c.Next() + `">Next</A>`))
}
