package render

import "io"

func StartHtml(w io.Writer) {
	w.Write([]byte(`<HTML><HEAD><STYLE TYPE="text/css" REL="stylesheet">body { font-family: "Comic Sans MS", "Comic Sans", cursive; }</STYLE></HEAD><BODY>`))
}
func EndHtml(w io.Writer) {
	w.Write([]byte(`</BODY><HTML>`))
}
func StartTable(w io.Writer) {
	w.Write([]byte(`<TABLE CELLPADDING="10" BORDER="1">`))
}
func EndTable(w io.Writer) {
	w.Write([]byte(`</TABLE>`))
}
