package render

import "io"

func StartHtml(w io.Writer) {
	w.Write([]byte(`<html><body>`))
}
func EndHtml(w io.Writer) {
	w.Write([]byte(`</body><html>`))
}
func StartTable(w io.Writer) {
	w.Write([]byte(`<table cellpadding="10" border="1">`))
}
func EndTable(w io.Writer) {
	w.Write([]byte(`</table>`))
}
