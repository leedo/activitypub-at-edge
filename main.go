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
		p := newProxy(w, r)
		remoteUrl, err := p.RemoteUrl()
		if err != nil {
			p.handleError(fsthttp.StatusBadRequest, err.Error())
			return
		}

		o, err := p.c.GetObject(ctx, remoteUrl)
		if err != nil {
			p.handleError(fsthttp.StatusBadRequest, err.Error())
			return
		}

		switch o.Type() {
		case "Person":
			p.renderPerson(ctx, o.ToPerson())
		case "Note":
			p.renderNote(ctx, o.ToNote())
		default:
			p.handleError(fsthttp.StatusBadRequest, "unknown object type")
		}
	})
}

func newProxy(w fsthttp.ResponseWriter, r *fsthttp.Request) *proxy {
	return &proxy{activitypub.NewClient(), w, r}
}

func (p *proxy) handleError(status int, msg string) {
	p.w.WriteHeader(status)
	p.w.Write([]byte(msg))
}

func (p *proxy) RemoteUrl() (string, error) {
	if p.r.Method != "GET" {
		return "", fmt.Errorf("This method is not allowed")
	}
	if p.r.URL.Path == "/favicon.ico" {
		return "", fmt.Errorf("Not Found")
	}

	u, err := url.Parse(p.r.URL.Path[1:]) // strip leading slash
	if err != nil {
		return "", fmt.Errorf("Invalid URL")
	}

	return u.String(), nil
}

func (p *proxy) renderNote(ctx context.Context, n *activitypub.Note) {
	person, err := p.c.GetPerson(ctx, n.AttributedTo())
	if err != nil {
		p.handleError(fsthttp.StatusBadRequest, err.Error())
		return
	}

	p.w.Header().Add("Content-Type", htmlType)
	p.w.WriteHeader(fsthttp.StatusOK)

	render.StartHtml(p.w)
	render.StartTable(p.w)
	render.Note(p.w, person, n)
	render.EndTable(p.w)
	render.EndHtml(p.w)
}

func (p *proxy) renderPerson(ctx context.Context, person *activitypub.Person) {
	col, err := p.c.GetCollection(ctx, person.Outbox())
	if err != nil {
		p.handleError(fsthttp.StatusBadRequest, err.Error())
		return
	}

	col, err = p.c.GetCollection(ctx, col.First())
	if err != nil {
		p.handleError(fsthttp.StatusBadRequest, err.Error())
		return
	}

	p.w.Header().Add("Content-Type", htmlType)
	p.w.WriteHeader(fsthttp.StatusOK)

	render.StartHtml(p.w)
	render.Person(p.w, person)
	render.StartTable(p.w)

	for _, item := range col.CollectionItems() {
		switch item.Type() {
		case "Create":
			switch o := item.Object(); o.Type() {
			case "Note":
				render.Note(p.w, person, o.ToNote())
			default:
				render.Unknown(p.w, o)
			}
		case "Announce":
			obj := item.Object()
			render.Unknown(p.w, obj)
		}
	}

	render.EndTable(p.w)
	render.EndHtml(p.w)
}
