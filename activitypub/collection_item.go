package activitypub

import "github.com/valyala/fastjson"

func (c *CollectionItem) Actor() string { return string(c.json.GetStringBytes("actor")) }
func (c *CollectionItem) Type() string  { return string(c.json.GetStringBytes("type")) }

func (c *CollectionItem) Object() *Object {
	if c.json.GetObject("object") != nil {
		v := c.json.Get("object")
		return &Object{v, v.GetStringBytes("id")}
	}

	return &Object{&fastjson.Value{}, c.json.GetStringBytes("object")}
}
