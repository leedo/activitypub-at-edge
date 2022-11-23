package render

import "io"

func Error(w io.Writer, err error) {
	w.Write([]byte(err.Error()))
}
