package render

import (
	"io"

	"github.com/leedo/activitypub-at-edge/user"
)

func Subscriptions(w io.Writer, s *user.Settings) {
	w.Write([]byte(`<H2>Subscriptions:</H2>`))
	if len(s.Subscriptions) > 0 {
		w.Write([]byte(`<UL>`))
		for _, f := range s.Subscriptions {
			w.Write([]byte(`<LI><A HREF="/` + f + `">` + f + `</A></LI>`))
		}
		w.Write([]byte(`</UL>`))
	} else {
		w.Write([]byte(`<P><I>No subscriptions</I></P>`))
	}
}
