package activitypub

import "github.com/valyala/fastjson"

func (o *Question) Type() string         { return string(o.json.GetStringBytes("type")) }
func (o *Question) ID() string           { return string(o.json.GetStringBytes("id")) }
func (o *Question) URL() string          { return string(o.json.GetStringBytes("url")) }
func (o *Question) Content() []byte      { return o.json.GetStringBytes("content") }
func (o *Question) Published() string    { return string(o.json.GetStringBytes("published")) }
func (o *Question) AttributedTo() string { return string(o.json.GetStringBytes("attributedTo")) }
func (o *Question) Loaded() bool         { return o.json != nil }

func (o *Question) InReplyTo() *Object {
	if v := o.json.Get("inReplyTo"); v != nil && v.Type() != fastjson.TypeNull {
		return &Object{v}
	}
	return nil
}
