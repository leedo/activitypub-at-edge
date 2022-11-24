package render

import (
	"io"
	"os"
)

func Footer(w io.Writer) {
	w.Write([]byte(`<P><I>Version ` + os.Getenv("FASTLY_SERVICE_VERSION") + `. Powered by Fastly Compute@Edge</I></P>`))
}
