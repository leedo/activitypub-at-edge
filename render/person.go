package render

import (
	"bufio"
	"io"

	"github.com/leedo/activitypub-at-edge/activitypub"
)

func Person(w io.Writer, p *activitypub.Person) {
	b := bufio.NewWriter(w)
	b.WriteString(`<A HREF="/` + p.ID() + `"><IMG SRC="` + p.Image().URL + `" HEIGHT="200"></A>`)
	b.WriteString(`<H1><A HREF="/` + p.ID() + `">` + p.Name() + `</A>'s Page</H1>`)
	b.Flush()
}
