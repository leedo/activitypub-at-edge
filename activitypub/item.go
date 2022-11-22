package activitypub

import "github.com/valyala/fastjson"

func (o *Item) Actor() string { return string(o.json.GetStringBytes("actor")) }

func (i *Item) Object() *Object {
	o := &Object{}
	if i.json.GetObject("object") != nil {
		o.json = i.json.Get("object")
		o.ID = o.json.GetStringBytes("id")
	} else {
		o.json = &fastjson.Value{}
		o.ID = i.json.GetStringBytes("object")
	}
	return o
}
