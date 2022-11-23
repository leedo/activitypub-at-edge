package activitypub

import "github.com/valyala/fastjson"

func (c *CollectionItem) Actor() string { return string(c.json.GetStringBytes("actor")) }
func (c *CollectionItem) Type() string  { return string(c.json.GetStringBytes("type")) }

func (c *CollectionItem) Get(q string) *fastjson.Value {
	return c.json.Get(q)
}
