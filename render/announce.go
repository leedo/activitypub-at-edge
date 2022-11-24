package render

import (
	"bufio"
	"io"

	"github.com/leedo/activitypub-at-edge/activitypub"
)

func Announce(w io.Writer, o *activitypub.Object) {
	b := bufio.NewWriter(w)
	b.WriteString(`<TR><TD></TD><TD><A HREF="` + o.ID() + `">` + o.ID() + `</A></TD></TR>`)
	b.Flush()
}
