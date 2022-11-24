package render

import (
	"fmt"
	"io"
)

func Error(w io.Writer, err string) {
	w.Write([]byte(fmt.Sprintf(`<tr><td></td><td>%s</td></tr>`, err)))
}
