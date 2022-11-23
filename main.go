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

type proxy struct {
	c *activitypub.Client
	w fsthttp.ResponseWriter
	r *fsthttp.Request
}

func main() {
	fsthttp.ServeFunc(func(ctx context.Context, w fsthttp.ResponseWriter, r *fsthttp.Request) {
		i := newProxy(w, r)
		remoteUrl, err := i.RemoteUrl()
		if err != nil {
			i.handleError(fsthttp.StatusBadRequest, err.Error())
			return
		}

		o, err := i.c.GetObject(ctx, remoteUrl)
		if err != nil {
			i.handleError(fsthttp.StatusBadRequest, err.Error())
			return
		}

		switch o.Type() {
		case "Person":
			i.renderPerson(ctx, o.ToPerson())
		case "Note":
			i.renderNote(ctx, o.ToNote())
		default:
			i.handleError(fsthttp.StatusBadRequest, "unknown object type")
		}
	})
}

func newProxy(w fsthttp.ResponseWriter, r *fsthttp.Request) *proxy {
	return &proxy{activitypub.NewClient(), w, r}
}

func (i *proxy) handleError(status int, msg string) {
	i.w.WriteHeader(status)
	i.w.Write([]byte(msg))
}

func (i *proxy) RemoteUrl() (string, error) {
	if i.r.Method != "GET" {
		return "", fmt.Errorf("This method is not allowed")
	}
	if i.r.URL.Path == "/favicon.ico" {
		return "", fmt.Errorf("Not Found")
	}

	u, err := url.Parse(i.r.URL.Path[1:]) // strip leading slash
	if err != nil {
		return "", fmt.Errorf("Invalid URL")
	}

	return u.String(), nil
}

func (i *proxy) renderNote(ctx context.Context, n *activitypub.Note) {
	p, err := i.c.GetPerson(ctx, n.AttributedTo())
	if err != nil {
		i.handleError(fsthttp.StatusBadRequest, err.Error())
		return
	}

	i.w.Header().Add("Content-Type", htmlType)
	i.w.WriteHeader(fsthttp.StatusOK)

	render.StartHtml(i.w)
	render.StartTable(i.w)
	render.Note(i.w, p, n)
	render.EndTable(i.w)
	render.EndHtml(i.w)
}

func (i *proxy) renderPerson(ctx context.Context, p *activitypub.Person) {
	col, err := i.c.GetCollection(ctx, p.Outbox())
	if err != nil {
		i.handleError(fsthttp.StatusBadRequest, err.Error())
		return
	}

	col, err = i.c.GetCollection(ctx, col.First())
	if err != nil {
		i.handleError(fsthttp.StatusBadRequest, err.Error())
		return
	}

	i.w.Header().Add("Content-Type", htmlType)
	i.w.WriteHeader(fsthttp.StatusOK)

	render.StartHtml(i.w)
	render.Person(i.w, p)
	render.StartTable(i.w)

	for _, item := range col.CollectionItems() {
		switch item.Type() {
		case "Create":
			switch o := item.Object(); o.Type() {
			case "Note":
				render.Note(i.w, p, o.ToNote())
			default:
				render.Unknown(i.w, o)
			}
		case "Announce":
			obj := item.Object()
			render.Unknown(i.w, obj)
		}
	}

	render.EndTable(i.w)
	render.EndHtml(i.w)
}
