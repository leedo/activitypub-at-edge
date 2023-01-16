package render

import (
	"bufio"
	"io"

	"github.com/leedo/activitypub-at-edge/activitypub"
	"github.com/leedo/activitypub-at-edge/user"
)

func PersonHeader(w io.Writer, p *activitypub.Person, s *user.Settings) {
	b := bufio.NewWriter(w)
	b.WriteString(`<A HREF="/` + p.URL() + `"><IMG SRC="` + p.Image().URL + `" HEIGHT="200"></A>`)
	b.WriteString(`<H1><A HREF="/` + p.URL() + `">` + p.Name() + `</A>'s Page</H1>`)
	b.WriteString(p.Summary())
	b.WriteString(`<FORM ACTION="/subscriptions" METHOD="POST">`)
	b.WriteString(`<INPUT TYPE="hidden" NAME="url" VALUE="` + p.URL() + `">`)
	if s.IsSubscribed(p.URL()) {
		b.WriteString(`<INPUT TYPE="hidden" NAME="action" VALUE="remove">`)
		b.WriteString(`<INPUT TYPE="submit" VALUE="Unsubscribe">`)

	} else {
		b.WriteString(`<INPUT TYPE="hidden" NAME="action" VALUE="add">`)
		b.WriteString(`<INPUT TYPE="submit" VALUE="Subscribe">`)
	}
	b.WriteString(`</FORM>`)
	b.Flush()
}
