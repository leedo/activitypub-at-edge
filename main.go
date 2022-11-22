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

type instance struct {
	c *activitypub.Client
	w fsthttp.ResponseWriter
	r *fsthttp.Request
}

func main() {
	fsthttp.ServeFunc(func(ctx context.Context, w fsthttp.ResponseWriter, r *fsthttp.Request) {
		i := newInstance(w, r)
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

func (i *instance) handleError(status int, msg string) {
	i.w.WriteHeader(status)
	i.w.Write([]byte(msg))
}

func (i *instance) RemoteUrl() (string, error) {
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

func newInstance(w fsthttp.ResponseWriter, r *fsthttp.Request) *instance {
	return &instance{activitypub.NewClient(), w, r}
}

func (i *instance) renderNote(ctx context.Context, n *activitypub.Note) {
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

func (i *instance) renderPerson(ctx context.Context, p *activitypub.Person) {
	o, err := i.c.GetCollection(ctx, p.Outbox())
	if err != nil {
		i.handleError(fsthttp.StatusBadRequest, err.Error())
		return
	}

	o, err = i.c.GetCollection(ctx, o.First())
	if err != nil {
		i.handleError(fsthttp.StatusBadRequest, err.Error())
		return
	}

	i.w.Header().Add("Content-Type", htmlType)
	i.w.WriteHeader(fsthttp.StatusOK)

	render.StartHtml(i.w)
	render.Person(i.w, p)
	render.StartTable(i.w)

	for _, item := range o.CollectionItems() {
		switch item.Type() {
		case "Create":
			obj := item.Object()
			switch obj.Type() {
			case "Note":
				render.Note(i.w, p, obj.ToNote())
			default:
				render.Unknown(i.w, obj)
			}
		case "Announce":
			obj := item.Object()
			render.Unknown(i.w, obj)
		}
	}

	render.EndTable(i.w)
	render.EndHtml(i.w)
}
