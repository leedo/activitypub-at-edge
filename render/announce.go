package render

import (
	"bufio"
	"io"

	"github.com/leedo/activitypub-at-edge/activitypub"
)

func Announce(w io.Writer, o *activitypub.Object) {
	b := bufio.NewWriter(w)
	remoteUrl := string(o.ID)
	b.WriteString(`<TR><TD></TD><TD><A HREF="` + remoteUrl + `">` + remoteUrl + `</A></TD></TR>`)
	b.Flush()
}
