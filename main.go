package main

import (
	"context"
	"fmt"
	"net/url"

	"github.com/fastly/compute-sdk-go/fsthttp"
	"github.com/leedo/activitypub-at-edge/activitypub"
	"github.com/leedo/activitypub-at-edge/render"
)

const htmlType = `text/html; charset="UTF-8"`

func main() {
	fsthttp.ServeFunc(func(ctx context.Context, w fsthttp.ResponseWriter, r *fsthttp.Request) {
		if r.Method != "GET" {
			w.WriteHeader(fsthttp.StatusMethodNotAllowed)
			fmt.Fprintf(w, "this method is not allowed")
			return
		}
		if r.URL.Path == "/favicon.ico" {
			w.WriteHeader(fsthttp.StatusNotFound)
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

		o, err := c.GetObject(ctx, remoteUrl.String())
		if err != nil {
			w.WriteHeader(fsthttp.StatusBadRequest)
			return
		}

		switch o.Type() {
		case "Person":
			renderPerson(ctx, w, c, o.ToPerson())
		case "Note":
			renderNote(ctx, w, c, o.ToNote())
		default:
			w.WriteHeader(fsthttp.StatusBadRequest)
			return
		}
	})
}

func renderNote(ctx context.Context, w fsthttp.ResponseWriter, c *activitypub.Client, n *activitypub.Note) {
	p, err := c.GetPerson(ctx, n.AttributedTo())
	if err != nil {
		w.WriteHeader(fsthttp.StatusBadRequest)
		fmt.Fprintf(w, "error fetching outbox: %s", err)
		return
	}

	w.Header().Add("Content-Type", htmlType)
	w.WriteHeader(fsthttp.StatusOK)

	render.StartHtml(w)
	render.StartTable(w)
	render.Note(w, p, n)
	render.EndTable(w)
	render.EndHtml(w)
}

func renderPerson(ctx context.Context, w fsthttp.ResponseWriter, c *activitypub.Client, p *activitypub.Person) {
	o, err := c.GetCollection(ctx, p.Outbox())
	if err != nil {
		w.WriteHeader(fsthttp.StatusBadRequest)
		fmt.Fprintf(w, "error fetching outbox: %s", err)
		return
	}

	o, err = c.GetCollection(ctx, o.First())
	if err != nil {
		w.WriteHeader(fsthttp.StatusBadRequest)
		fmt.Fprintf(w, "error fetching outbox: %s", err)
		return
	}

	w.Header().Add("Content-Type", htmlType)
	w.WriteHeader(fsthttp.StatusOK)

	render.StartHtml(w)
	render.Person(w, p)
	render.StartTable(w)

	for _, i := range o.CollectionItems() {
		switch i.Type() {
		case "Create":
			obj := i.Object()
			switch obj.Type() {
			case "Note":
				render.Note(w, p, obj.ToNote())
			default:
				render.Unknown(w, obj)
			}
		case "Announce":
			obj := i.Object()
			render.Unknown(w, obj)
		}
	}

	render.EndTable(w)
	render.EndHtml(w)
}
