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
	backends map[string]backend
	cache    map[string]*Object
}

type backend struct {
	name string
}

func (c *Client) AddBackend(name string) {
	c.backends[name] = backend{name}
}

func (c *Client) GetObject(ctx context.Context, remoteUrl string) (*Object, error) {
	if o, ok := c.cache[remoteUrl]; ok {
		return o, nil
	}

	v, err := c.request(ctx, "GET", remoteUrl, nil)
	if err != nil {
		return nil, err
	}
	o := NewObject(v)
	c.cache[remoteUrl] = o
	return o, nil
}

func (c *Client) GetCollection(ctx context.Context, remoteUrl string) (*Collection, error) {
	o, err := c.GetObject(ctx, remoteUrl)
	if err != nil {
		return nil, err
	}
	return o.ToCollection(), nil
}

func (c *Client) GetPerson(ctx context.Context, remoteUrl string) (*Person, error) {
	o, err := c.GetObject(ctx, remoteUrl)
	if err != nil {
		return nil, err
	}
	return o.ToPerson(), nil
}

func (c *Client) request(ctx context.Context, method string, remoteUrl string, body io.Reader) (*fastjson.Value, error) {
	req, err := fsthttp.NewRequest(method, remoteUrl, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", activityJSON)
	req.CacheOptions.TTL = 900

	b, ok := c.backends[req.URL.Host]
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
		backends: make(map[string]backend, 0),
		cache:    make(map[string]*Object, 0),
	}
}
