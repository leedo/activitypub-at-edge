package render

import "io"

func Login(w io.Writer) {
	StartHtml(w, nil)
	w.Write([]byte(`<A HREF="/gh_login">Login with Github</A>`))
	EndHtml(w)
}
