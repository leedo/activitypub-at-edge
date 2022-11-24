package render

import (
	"bufio"
	"io"

	"github.com/leedo/activitypub-at-edge/activitypub"
)

func Announce(w io.Writer, a *activitypub.Person) {
	b := bufio.NewWriter(w)
	b.WriteString(`<TR><TD COLSPAN="2"><A HREF="` + a.URL() + `">` + a.Name() + `</A> announced:</TD></TR>`)
	b.Flush()
}
