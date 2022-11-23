package render

import (
	"io"
)

func Unknown(w io.Writer, t string) {
	w.Write([]byte(`<tr><td></td><td>Unable to render unknown type ` + t + `</td></tr>`))
}
