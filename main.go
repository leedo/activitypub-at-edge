package main

import (
	"context"
	"fmt"
	"net/url"

	"github.com/fastly/compute-sdk-go/fsthttp"
	"github.com/leedo/activitypub-at-edge/activitypub"
)

func main() {
	fsthttp.ServeFunc(func(ctx context.Context, w fsthttp.ResponseWriter, r *fsthttp.Request) {
		if r.Method != "GET" {
			w.WriteHeader(fsthttp.StatusMethodNotAllowed)
			fmt.Fprintf(w, "this method is not allowed")
			return
		}

		remoteUrl, err := url.Parse(r.URL.Path[1:]) // strip leading slash
		if err != nil {
			w.WriteHeader(fsthttp.StatusBadRequest)
			fmt.Fprintf(w, "invalid URL: %s", err)
			return
		}

		c := activitypub.NewClient()
		c.AddBackend(remoteUrl.Host)

		p, err := c.GetPerson(ctx, remoteUrl.String())
		if err != nil {
			w.WriteHeader(fsthttp.StatusBadRequest)
			fmt.Fprintf(w, "error fetching person: %s", err)
			return
		}

		o, err := c.GetOutbox(ctx, p.Outbox())
		if err != nil {
			w.WriteHeader(fsthttp.StatusBadRequest)
			fmt.Fprintf(w, "error fetching outbox: %s", err)
			return
		}

		o, err = c.GetOutbox(ctx, o.First())
		if err != nil {
			w.WriteHeader(fsthttp.StatusBadRequest)
			fmt.Fprintf(w, "error fetching outbox: %s", err)
			return
		}

		w.Header().Add("Content-Type", "text/html; charset=\"UTF-8\"")
		w.WriteHeader(fsthttp.StatusOK)

		w.Write([]byte("<html><body><table border=\"1\">"))
		for _, i := range o.Items() {
			obj := i.Object()
			w.Write([]byte("<tr><td>"))
			person, err := c.GetPerson(ctx, obj.AttributedTo())
			if err == nil {
				if i := person.Icon(); i.URL != "" {
					w.Write([]byte("<img width=\"100\" src=\"" + i.URL + "\"><br>"))
				}
				w.Write([]byte(person.Name()))
			}

			w.Write([]byte("</td><td>"))
			if err := c.LoadObject(obj); err != nil {
				fmt.Fprintf(w, "%s: %s", obj.ID, err)
			} else {
				w.Write(obj.Content())
			}
			w.Write([]byte("</td></tr>"))
		}
		w.Write([]byte("</table></body></html>"))
	})
}
