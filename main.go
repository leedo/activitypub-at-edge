package main

import (
	"context"
	"fmt"
	"io"
	"net/url"

	"activitypub"

	"activitypub/activitypub"

	"github.com/fastly/compute-sdk-go/fsthttp"
	"github.com/valyala/fastjson"
)

const activityJSON = "application/activity+json"

type server struct {
	backends    map[string]backend
	personCache map[string]*activitypub.Person
}

type backend struct {
	name string
}

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

		s := NewServer()
		s.addBackend(remoteUrl.Host)

		p, err := s.GetPerson(ctx, remoteUrl.String())
		if err != nil {
			w.WriteHeader(fsthttp.StatusBadRequest)
			fmt.Fprintf(w, "error fetching person: %s", err)
			return
		}

		o, err := s.GetOutbox(ctx, p.Outbox())
		if err != nil {
			w.WriteHeader(fsthttp.StatusBadRequest)
			fmt.Fprintf(w, "error fetching outbox: %s", err)
			return
		}

		o, err = s.GetOutbox(ctx, o.First())
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
			person, err := s.GetPerson(ctx, obj.AttributedTo())
			if err == nil {
				if i := person.Icon(); i.url != "" {
					w.Write([]byte("<img width=\"100\" src=\"" + i.url + "\"><br>"))
				}
				w.Write([]byte(person.Name()))
			}

			w.Write([]byte("</td><td>"))
			if err := s.LoadObject(obj); err != nil {
				fmt.Fprintf(w, "%s: %s", obj.id, err)
			} else {
				w.Write(obj.Content())
			}
			w.Write([]byte("</td></tr>"))
		}
		w.Write([]byte("</table></body></html>"))
	})
}

func NewServer() *server {
	return &server{
		backends:    make(map[string]backend, 0),
		personCache: make(map[string]*activitypub.Person, 0),
	}
}

func (s *server) addBackend(name string) {
	s.backends[name] = backend{name}
}

func (s *server) LoadObject(o *activitypub.Object) error {
	if o.json != nil {
		return nil
	}

	return fmt.Errorf("remote server not supported")
}

func (s *server) GetObject(ctx context.Context, remoteUrl string) (*activitypub.Object, error) {
	v, err := s.request(ctx, "GET", remoteUrl, nil)
	if err != nil {
		return nil, err
	}
	return &activitypub.Object{v, v.GetStringBytes("id")}, nil
}

func (s *server) GetOutbox(ctx context.Context, remoteUrl string) (*activitypub.Outbox, error) {
	v, err := s.request(ctx, "GET", remoteUrl, nil)
	if err != nil {
		return nil, err
	}
	return &activitypub.Outbox{v, v.GetStringBytes("id")}, nil
}

func (s *server) GetPerson(ctx context.Context, remoteUrl string) (*activitypub.Person, error) {
	if p, ok := s.personCache[remoteUrl]; ok {
		return p, nil
	}

	v, err := s.request(ctx, "GET", remoteUrl, nil)
	if err != nil {
		return nil, err
	}
	p := &activitypub.Person{v, v.GetStringBytes("id")}
	s.personCache[remoteUrl] = p
	return p, nil
}

func (s *server) request(ctx context.Context, method string, remoteUrl string, body io.Reader) (*fastjson.Value, error) {
	req, err := fsthttp.NewRequest(method, remoteUrl, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", activityJSON)
	req.CacheOptions.TTL = 900

	b, ok := s.backends[req.URL.Host]
	if !ok {
		return nil, fmt.Errorf("unknown backend %q", req.URL.Host)
	}

	resp, err := req.Send(ctx, b.name)
	if err != nil {
		return nil, err
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var p fastjson.Parser
	v, err := p.ParseBytes(respBody)
	if err != nil {
		return nil, err
	}

	return v, nil
}
