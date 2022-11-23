package proxy

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/fastly/compute-sdk-go/fsthttp"
	"github.com/leedo/activitypub-at-edge/activitypub"
	"github.com/leedo/activitypub-at-edge/render"
)

const htmlType = `text/html; charset="UTF-8"`

type proxy struct {
	c *activitypub.Client
	w fsthttp.ResponseWriter
}

func NewProxy(w fsthttp.ResponseWriter) *proxy {
	return &proxy{activitypub.NewClient(), w}
}

func (p *proxy) debug(msg string) {
	fmt.Fprintf(os.Stderr, "%s\n", msg)
}

func (p *proxy) ErrorHandler(status int, msg string) {
	p.w.WriteHeader(status)
	p.w.Write([]byte(msg))
}

func (p *proxy) remoteUrl(r *fsthttp.Request) (string, error) {
	if r.Method != "GET" {
		return "", fmt.Errorf("This method is not allowed")
	}
	if r.URL.Path == "/favicon.ico" {
		return "", fmt.Errorf("Not Found")
	}

	u, err := url.Parse(r.URL.RequestURI()[1:]) // strip leading slash
	if err != nil {
		return "", fmt.Errorf("Invalid URL")
	}

	return u.String(), nil
}

func (p *proxy) NoteHandler(ctx context.Context, n *activitypub.Note) {
	person, err := p.c.GetPerson(ctx, n.AttributedTo())
	if err != nil {
		p.ErrorHandler(fsthttp.StatusBadRequest, err.Error())
		return
	}

	p.w.Header().Add("Content-Type", htmlType)
	p.w.WriteHeader(fsthttp.StatusOK)

	render.StartHtml(p.w)
	render.StartTable(p.w)
	render.Note(p.w, person, n)
	render.EndTable(p.w)
	render.Footer(p.w)
	render.EndHtml(p.w)
}

func (p *proxy) CollectionHandler(ctx context.Context, col *activitypub.Collection) {
	p.w.Header().Add("Content-Type", htmlType)
	p.w.WriteHeader(fsthttp.StatusOK)

	render.StartHtml(p.w)
	render.Pagination(p.w, col)
	render.StartTable(p.w)

	for _, item := range col.CollectionItems() {
		p.renderCollectionItem(ctx, item)
	}

	render.EndTable(p.w)
	render.Pagination(p.w, col)
	render.Footer(p.w)
	render.EndHtml(p.w)
}

func (p *proxy) PersonHandler(ctx context.Context, person *activitypub.Person) {
	col, err := p.c.GetCollection(ctx, person.Outbox())
	if err != nil {
		p.ErrorHandler(fsthttp.StatusBadRequest, err.Error())
		return
	}

	col, err = p.c.GetCollection(ctx, col.First())
	if err != nil {
		p.ErrorHandler(fsthttp.StatusBadRequest, err.Error())
		return
	}

	p.w.Header().Add("Content-Type", htmlType)
	p.w.WriteHeader(fsthttp.StatusOK)

	render.StartHtml(p.w)
	render.Person(p.w, person)
	render.Pagination(p.w, col)
	render.StartTable(p.w)

	for _, item := range col.CollectionItems() {
		p.renderCollectionItem(ctx, item)
	}

	render.EndTable(p.w)
	render.Pagination(p.w, col)
	render.Footer(p.w)
	render.EndHtml(p.w)
}

func (p *proxy) renderCollectionItem(ctx context.Context, item *activitypub.CollectionItem) {
	switch item.Type() {
	case "Create":
		o, err := p.c.NewObject(ctx, item.Get("object"))
		if err != nil {
			render.Error(p.w, err)
			return
		}
		switch o.Type() {
		case "Note":
			note := o.ToNote()
			person, err := p.c.GetPerson(ctx, note.AttributedTo())
			if err != nil {
				render.Error(p.w, err)
			} else {
				render.Note(p.w, person, note)
			}
		default:
			render.Unknown(p.w, o.Type())
		}
	case "Announce":
		render.Unknown(p.w, item.Type())
	}
}

func (p *proxy) GenericRequestHandler(ctx context.Context, r *fsthttp.Request) {
	remoteUrl, err := p.remoteUrl(r)
	if err != nil {
		p.ErrorHandler(fsthttp.StatusBadRequest, err.Error())
		return
	}

	o, err := p.c.GetObject(ctx, remoteUrl)
	if err != nil {
		p.ErrorHandler(fsthttp.StatusBadRequest, err.Error())
		return
	}

	switch o.Type() {
	case activitypub.PersonType:
		p.PersonHandler(ctx, o.ToPerson())
	case activitypub.NoteType:
		p.NoteHandler(ctx, o.ToNote())
	case activitypub.OrderedCollectionPageType:
		p.CollectionHandler(ctx, o.ToCollection())
	default:
		p.ErrorHandler(fsthttp.StatusBadRequest, "unknown object type "+o.Type())
	}
}
