package render

import (
	"bufio"
	"io"

	"github.com/leedo/activitypub-at-edge/activitypub"
)

func Announce(w io.Writer, o *activitypub.Object) {
	b := bufio.NewWriter(w)
	remoteUrl := string(o.ID)
	b.WriteString(`<tr><td></td><td><a href="` + remoteUrl + `">` + remoteUrl + `</a></td></tr>`)
	b.Flush()
}
