package activitypub

import (
	"context"
	"fmt"
	"io"

	"github.com/fastly/compute-sdk-go/fsthttp"
	"github.com/valyala/fastjson"
)

const activityJSON = "application/activity+json"

type Client struct {
	backends    map[string]backend
	personCache map[string]*Person
}

type backend struct {
	name string
}

func (s *Client) AddBackend(name string) {
	s.backends[name] = backend{name}
}

func (s *Client) GetObject(ctx context.Context, i *Item) (*Object, error) {
	o := i.Object()
	if o.Loaded() {
		return o, nil
	}

	v, err := s.request(ctx, "GET", string(o.ID), nil)
	if err != nil {
		return nil, err
	}
	return NewObject(v), nil
}

func (s *Client) GetOutbox(ctx context.Context, remoteUrl string) (*Outbox, error) {
	v, err := s.request(ctx, "GET", remoteUrl, nil)
	if err != nil {
		return nil, err
	}
	return NewOutbox(v), nil
}

func (s *Client) GetPerson(ctx context.Context, remoteUrl string) (*Person, error) {
	if p, ok := s.personCache[remoteUrl]; ok {
		return p, nil
	}

	v, err := s.request(ctx, "GET", remoteUrl, nil)
	if err != nil {
		return nil, err
	}
	p := NewPerson(v)
	s.personCache[remoteUrl] = p
	return p, nil
}

func (s *Client) request(ctx context.Context, method string, remoteUrl string, body io.Reader) (*fastjson.Value, error) {
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

func NewClient() *Client {
	return &Client{
		backends:    make(map[string]backend, 0),
		personCache: make(map[string]*Person, 0),
	}
}
