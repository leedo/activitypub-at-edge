package server

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/fastly/compute-sdk-go/fsthttp"
	"github.com/leedo/activitypub-at-edge/activitypub"
	"github.com/leedo/activitypub-at-edge/oauth"
	"github.com/leedo/activitypub-at-edge/render"
)

const htmlType = `text/html; charset="UTF-8"`

type Server struct {
	a *oauth.OAuth
	c *activitypub.Client
	w fsthttp.ResponseWriter
}

func NewServer(w fsthttp.ResponseWriter, a *oauth.OAuth) *Server {
	return &Server{a, activitypub.NewClient(), w}
}

func (s *Server) debug(msg string) {
	fmt.Fprintf(os.Stderr, "%s\n", msg)
}

func (s *Server) ErrorPage(status int, msg string) {
	s.w.WriteHeader(status)
	render.StartHtml(s.w, s.a.User)
	render.StartTable(s.w)
	render.Error(s.w, msg)
	render.EndTable(s.w)
	render.Footer(s.w)
	render.EndHtml(s.w)
}

func (s *Server) remoteUrl(r *fsthttp.Request) (string, error) {
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

func (s *Server) NotePage(ctx context.Context, o *activitypub.Object) {
	s.w.Header().Add("Content-Type", htmlType)
	s.w.WriteHeader(fsthttp.StatusOK)

	render.StartHtml(s.w, s.a.User)
	render.StartTable(s.w)

	n := o.ToNote()

	if parent := n.InReplyTo(); parent != nil {
		s.renderObject(ctx, parent)
	}

	s.renderObject(ctx, o)

	render.EndTable(s.w)
	render.Footer(s.w)
	render.EndHtml(s.w)
}

func (s *Server) CollectionPage(ctx context.Context, col *activitypub.Collection) {
	s.w.Header().Add("Content-Type", htmlType)
	s.w.WriteHeader(fsthttp.StatusOK)

	render.StartHtml(s.w, s.a.User)
	s.renderCollection(ctx, col)
	render.Footer(s.w)
	render.EndHtml(s.w)
}

func (s *Server) renderCollection(ctx context.Context, col *activitypub.Collection) {
	render.Pagination(s.w, col)
	render.StartTable(s.w)

	for _, o := range col.CollectionItems() {
		s.renderObject(ctx, o)
	}

	render.EndTable(s.w)
	render.Pagination(s.w, col)
}

func (s *Server) PersonPage(ctx context.Context, person *activitypub.Person) {
	col, err := s.c.GetCollection(ctx, person.Outbox())
	if err != nil {
		s.ErrorPage(fsthttp.StatusBadGateway, err.Error())
		return
	}

	col, err = s.c.GetCollection(ctx, col.First())
	if err != nil {
		s.ErrorPage(fsthttp.StatusBadGateway, err.Error())
		return
	}

	s.w.Header().Add("Content-Type", htmlType)
	s.w.WriteHeader(fsthttp.StatusOK)

	render.StartHtml(s.w, s.a.User)
	render.PersonHeader(s.w, person)

	s.renderCollection(ctx, col)

	render.Footer(s.w)
	render.EndHtml(s.w)
}

func (s *Server) renderObject(ctx context.Context, o *activitypub.Object) {
	if err := s.c.LoadObject(ctx, o); err != nil {
		render.Error(s.w, err.Error())
		return
	}

	switch o.Type() {
	case activitypub.CreateType:
		s.renderObject(ctx, o.ToActivity().Object())

	case activitypub.AnnounceType:
		activity := o.ToActivity()
		person, err := s.c.GetPerson(ctx, activity.Actor())
		if err != nil {
			render.Error(s.w, err.Error())
		} else {
			render.Announce(s.w, person)
			s.renderObject(ctx, activity.Object())
		}

	case activitypub.NoteType:
		note := o.ToNote()
		person, err := s.c.GetPerson(ctx, note.AttributedTo())
		if err != nil {
			render.Error(s.w, err.Error())
		} else {
			render.Note(s.w, person, note)
		}

	default:
		render.Unknown(s.w, o.Type())
	}
}

func (s *Server) UserHandler(ctx context.Context, r *fsthttp.Request) {
	s.w.Header().Add("Content-Type", htmlType)
	s.w.WriteHeader(fsthttp.StatusOK)
	render.StartHtml(s.w, s.a.User)
	s.w.Write([]byte(s.a.User.Login))
	render.Footer(s.w)
	render.EndHtml(s.w)
}

func (s *Server) GenericRequestHandler(ctx context.Context, r *fsthttp.Request) {
	remoteUrl, err := s.remoteUrl(r)
	if err != nil {
		s.ErrorPage(fsthttp.StatusBadRequest, err.Error())
		return
	}

	o, err := s.c.GetObject(ctx, remoteUrl)
	if err != nil {
		s.ErrorPage(fsthttp.StatusBadRequest, err.Error())
		return
	}

	switch o.Type() {
	case activitypub.PersonType:
		s.PersonPage(ctx, o.ToPerson())
	case activitypub.NoteType:
		s.NotePage(ctx, o)
	case activitypub.OrderedCollectionPageType, activitypub.OrderedCollectionType:
		s.CollectionPage(ctx, o.ToCollection())
	default:
		s.ErrorPage(fsthttp.StatusNotFound, "unknown object type "+o.Type())
	}
}
