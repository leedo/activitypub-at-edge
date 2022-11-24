package render

import (
	"fmt"
	"io"
)

func Error(w io.Writer, err error) {
	w.Write([]byte(fmt.Sprintf(`<tr><td></td><td>Error: %s</td></tr>`, err.Error())))
}
