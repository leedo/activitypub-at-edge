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
	cache map[string]*Object
}

func (c *Client) GetObject(ctx context.Context, remoteUrl string) (*Object, error) {
	if o, ok := c.cache[remoteUrl]; ok {
		return o, nil
	}

	v, err := c.request(ctx, "GET", remoteUrl, nil)
	if err != nil {
		return nil, err
	}
	o := &Object{v}
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

	resp, err := req.Send(ctx, req.URL.Host)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != fsthttp.StatusOK {
		return nil, fmt.Errorf("error fetching %s: %d", remoteUrl, resp.StatusCode)
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

func (c *Client) LoadObject(ctx context.Context, o *Object) error {
	if !o.IsLoaded() {
		v, err := c.request(ctx, "GET", o.ID(), nil)
		if err != nil {
			return err
		}
		o.json = v
	}
	return nil
}

func (c *Client) NewObject(ctx context.Context, v *fastjson.Value) (*Object, error) {
	if v == nil {
		return nil, fmt.Errorf("unexpected nil JSON")
	}
	switch v.Type() {
	case fastjson.TypeString:
		return c.GetObject(ctx, string(v.GetStringBytes()))
	case fastjson.TypeObject:
		return &Object{v}, nil
	default:
		return nil, fmt.Errorf("unexpected JSON type for object %s", v.Type())
	}
}

func NewClient() *Client {
	return &Client{
		cache: make(map[string]*Object, 0),
	}
}
