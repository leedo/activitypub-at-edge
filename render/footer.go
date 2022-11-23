package render

import "io"

func Footer(w io.Writer) {
	w.Write([]byte(`<P><I>Powered by Fastly Compute@Edge</I></P>`))
}
