package render

import (
	"fmt"
	"io"

	"github.com/leedo/activitypub-at-edge/activitypub"
)

func Unknown(w io.Writer, o *activitypub.Object) {
	w.Write([]byte(fmt.Sprintf(`<tr><td></td><td>Unable to render unknown type %s %s</td></tr>`, o.Type(), string(o.ID))))
}
