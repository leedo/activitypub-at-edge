package render

import (
	"bufio"
	"io"

	"github.com/leedo/activitypub-at-edge/activitypub"
)

func Question(w io.Writer, p *activitypub.Person, n *activitypub.Question) {
	b := bufio.NewWriter(w)
	b.WriteString(`<TR><TD ALIGN="center" VALIGN="top" ROWSPAN="2">`)

	b.WriteString(`<A HREF="/` + p.URL() + `"><IMG WIDTH="100" SRC="` + p.Icon().URL + `"></A>`)
	b.WriteString(`<BR>`)
	b.WriteString(`<STRONG><a HREF="/` + p.URL() + `">` + p.Name() + `</a></STRONG>`)

	b.WriteString(`</TD><TD valign="top">`)
	b.Write(n.Content())
	b.WriteString(`</TD></TR><TR><TD>`)
	b.WriteString(`<A HREF="/` + n.URL() + `">` + n.Published() + `</A>`)
	b.WriteString(`</TD></TR>`)
	b.Flush()
}
