package render

import "io"

func Login(w io.Writer) {
	StartHtml(w, nil)
	w.Write([]byte(`<A HREF="/oauth_login">Login with Github</A>`))
	EndHtml(w)
}
