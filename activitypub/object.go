package activitypub

import "github.com/valyala/fastjson"

func (o *Object) Type() string         { return string(o.json.GetStringBytes("type")) }
func (o *Object) InReplyTo() string    { return string(o.json.GetStringBytes("inReplyTo")) }
func (o *Object) URL() string          { return string(o.json.GetStringBytes("url")) }
func (o *Object) Content() []byte      { return o.json.GetStringBytes("content") }
func (o *Object) Published() string    { return string(o.json.GetStringBytes("published")) }
func (o *Object) AttributedTo() string { return string(o.json.GetStringBytes("attributedTo")) }
func (o *Object) Loaded() bool         { return o.json != nil }

func (o *Object) Attachments() []*Attachment {
	var a []*Attachment
	for _, v := range o.json.GetArray("attachment") {
		a = append(a, &Attachment{
			URL:       string(v.GetStringBytes("url")),
			Type:      string(v.GetStringBytes("type")),
			MediaType: string(v.GetStringBytes("mediaType")),
			Width:     v.GetUint("width"),
			Height:    v.GetUint("height"),
		})
	}
	return a
}

func NewObject(v *fastjson.Value) *Object {
	return &Object{v, v.GetStringBytes("id")}
}
