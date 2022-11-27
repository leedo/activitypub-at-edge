package render

import (
	"io"

	"github.com/leedo/activitypub-at-edge/oauth"
)

func StartHtml(w io.Writer, u *oauth.User) {
	w.Write([]byte(`<HTML><HEAD><STYLE TYPE="text/css" REL="stylesheet">body { font-family: "Comic Sans MS", "Comic Sans", cursive; }</STYLE></HEAD><BODY>`))
	if u != nil {
		w.Write([]byte(`<P><IMG WIDTH=32 SRC="` + u.Avatar + `"> Welcome ` + u.Login + ` (<A HREF="/logout">Logout</A>)</P>`))
	}
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
