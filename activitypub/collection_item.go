package activitypub

import "github.com/valyala/fastjson"

func (o *CollectionItem) Actor() string { return string(o.json.GetStringBytes("actor")) }
func (o *CollectionItem) Type() string  { return string(o.json.GetStringBytes("type")) }

func (o *CollectionItem) Object() *Object {
	if o.json.GetObject("object") != nil {
		v := o.json.Get("object")
		return &Object{v, v.GetStringBytes("id")}
	}

	return &Object{&fastjson.Value{}, o.json.GetStringBytes("object")}
}
