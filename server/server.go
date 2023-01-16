package server

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"

	"github.com/fastly/compute-sdk-go/fsthttp"
	"github.com/leedo/activitypub-at-edge/activitypub"
	"github.com/leedo/activitypub-at-edge/render"
	"github.com/leedo/activitypub-at-edge/store"
	"github.com/leedo/activitypub-at-edge/user"
)

const htmlType = `text/html; charset="UTF-8"`

type Server struct {
	u *user.User
	c *activitypub.Client
	s *store.Store
}

func NewServer(u *user.User, s *store.Store) *Server {
	return &Server{u, activitypub.NewClient(), s}
}

func (s *Server) debug(msg string) {
	fmt.Fprintf(os.Stderr, "%s\n", msg)
}

func ErrorPage(status int, w fsthttp.ResponseWriter, msg string) {
	w.WriteHeader(status)
	render.StartHtml(w, nil)
	render.StartTable(w)
	render.Error(w, msg)
	render.EndTable(w)
	render.Footer(w)
	render.EndHtml(w)
}

func (s *Server) ErrorPage(status int, w fsthttp.ResponseWriter, msg string) {
	w.WriteHeader(status)
	render.StartHtml(w, s.u)
	render.StartTable(w)
	render.Error(w, msg)
	render.EndTable(w)
	render.Footer(w)
	render.EndHtml(w)
}

func (s *Server) remoteUrl(r *fsthttp.Request) (string, error) {
	if r.Method != "GET" {
		return "", fmt.Errorf("this method is not allowed")
	}
	if r.URL.Path == "/favicon.ico" {
		return "", fmt.Errorf("not found")
	}

	u, err := url.Parse(r.URL.RequestURI()[1:]) // strip leading slash
	if err != nil {
		return "", fmt.Errorf("invalid URL")
	}

	if u.Scheme != "https" {
		return "", fmt.Errorf("invalid URL")
	}

	return u.String(), nil
}

func (s *Server) NotePage(ctx context.Context, w fsthttp.ResponseWriter, o *activitypub.Object) {
	w.Header().Add("Content-Type", htmlType)
	w.WriteHeader(fsthttp.StatusOK)

	render.StartHtml(w, s.u)
	render.StartTable(w)

	n := o.ToNote()

	if parent := n.InReplyTo(); parent != nil {
		s.renderObject(ctx, w, parent)
	}

	s.renderObject(ctx, w, o)

	render.EndTable(w)
	render.Footer(w)
	render.EndHtml(w)
}

func (s *Server) CollectionPage(ctx context.Context, w fsthttp.ResponseWriter, col *activitypub.Collection) {
	w.Header().Add("Content-Type", htmlType)
	w.WriteHeader(fsthttp.StatusOK)

	render.StartHtml(w, s.u)
	s.renderCollection(ctx, w, col)
	render.Footer(w)
	render.EndHtml(w)
}

func (s *Server) renderCollection(ctx context.Context, w fsthttp.ResponseWriter, col *activitypub.Collection) {
	render.Pagination(w, col)
	render.StartTable(w)

	for _, o := range col.CollectionItems() {
		s.renderObject(ctx, w, o)
	}

	render.EndTable(w)
	render.Pagination(w, col)
}

func (s *Server) PersonPage(ctx context.Context, w fsthttp.ResponseWriter, person *activitypub.Person) {
	col, err := s.c.GetCollection(ctx, person.Outbox())
	if err != nil {
		s.ErrorPage(fsthttp.StatusInternalServerError, w, err.Error())
		return
	}

	col, err = s.c.GetCollection(ctx, col.First())
	if err != nil {
		s.ErrorPage(fsthttp.StatusInternalServerError, w, err.Error())
		return
	}

	settings, err := s.u.Settings()
	if err != nil {
		s.ErrorPage(fsthttp.StatusInternalServerError, w, err.Error())
		return
	}

	w.Header().Add("Content-Type", htmlType)
	w.WriteHeader(fsthttp.StatusOK)

	render.StartHtml(w, s.u)
	render.PersonHeader(w, person, settings)

	s.renderCollection(ctx, w, col)

	render.Footer(w)
	render.EndHtml(w)
}

func (s *Server) renderObject(ctx context.Context, w fsthttp.ResponseWriter, o *activitypub.Object) {
	if err := s.c.LoadObject(ctx, o); err != nil {
		render.Error(w, err.Error())
		return
	}

	switch o.Type() {
	case activitypub.CreateType:
		s.renderObject(ctx, w, o.ToActivity().Object())

	case activitypub.AnnounceType:
		activity := o.ToActivity()
		person, err := s.c.GetPerson(ctx, activity.Actor())
		if err != nil {
			render.Error(w, err.Error())
		} else {
			render.Announce(w, person)
			s.renderObject(ctx, w, activity.Object())
		}

	case activitypub.NoteType:
		note := o.ToNote()
		person, err := s.c.GetPerson(ctx, note.AttributedTo())
		if err != nil {
			render.Error(w, err.Error())
		} else {
			render.Note(w, person, note)
		}

	default:
		render.Unknown(w, o.Type())
	}
}

func (s *Server) UserHandler(ctx context.Context, w fsthttp.ResponseWriter, r *fsthttp.Request) {
	w.Header().Add("Content-Type", htmlType)
	w.WriteHeader(fsthttp.StatusOK)
	render.StartHtml(w, s.u)

	settings, err := s.u.Settings()
	if err != nil {
		render.Error(w, err.Error())
	} else {
		render.Subscriptions(w, settings)
	}

	render.Footer(w)
	render.EndHtml(w)
}

func (s *Server) SubscribeHandler(ctx context.Context, w fsthttp.ResponseWriter, r *fsthttp.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		ErrorPage(fsthttp.StatusInternalServerError, w, err.Error())
		return
	}

	q, err := url.ParseQuery(string(b))
	if err != nil {
		ErrorPage(fsthttp.StatusInternalServerError, w, err.Error())
		return
	}

	if !q.Has("url") {
		ErrorPage(fsthttp.StatusInternalServerError, w, "url is required")
		return
	}

	u := q.Get("url")
	if _, err := url.Parse(u); err != nil {
		ErrorPage(fsthttp.StatusInternalServerError, w, err.Error())
		return
	}

	if !q.Has("action") {
		ErrorPage(fsthttp.StatusInternalServerError, w, "action is required")
		return
	}

	switch q.Get("action") {
	case "add":
		if err := s.u.Subscribe(u); err != nil {
			ErrorPage(fsthttp.StatusInternalServerError, w, err.Error())
			return
		}
	case "remove":
		if err := s.u.Unsubscribe(u); err != nil {
			ErrorPage(fsthttp.StatusInternalServerError, w, err.Error())
			return
		}
	default:
		ErrorPage(fsthttp.StatusInternalServerError, w, "invalid action")
		return
	}

	w.Header().Set("Location", "/")
	w.WriteHeader(fsthttp.StatusFound)
	return

}

func (s *Server) GenericRequestHandler(ctx context.Context, w fsthttp.ResponseWriter, r *fsthttp.Request) {
	remoteUrl, err := s.remoteUrl(r)
	if err != nil {
		s.ErrorPage(fsthttp.StatusBadRequest, w, err.Error())
		return
	}

	o, err := s.c.GetObject(ctx, remoteUrl)
	if err != nil {
		s.ErrorPage(fsthttp.StatusBadRequest, w, err.Error())
		return
	}

	switch o.Type() {
	case activitypub.PersonType:
		s.PersonPage(ctx, w, o.ToPerson())
	case activitypub.NoteType:
		s.NotePage(ctx, w, o)
	case activitypub.OrderedCollectionPageType, activitypub.OrderedCollectionType:
		s.CollectionPage(ctx, w, o.ToCollection())
	default:
		s.ErrorPage(fsthttp.StatusNotFound, w, "unknown object type "+o.Type())
	}
}
