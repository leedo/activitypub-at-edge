package activitypub

import "github.com/valyala/fastjson"

func (o *Note) Type() string         { return string(o.json.GetStringBytes("type")) }
func (o *Note) ID() string           { return string(o.json.GetStringBytes("id")) }
func (o *Note) URL() string          { return string(o.json.GetStringBytes("url")) }
func (o *Note) Content() []byte      { return o.json.GetStringBytes("content") }
func (o *Note) Published() string    { return string(o.json.GetStringBytes("published")) }
func (o *Note) AttributedTo() string { return string(o.json.GetStringBytes("attributedTo")) }
func (o *Note) Loaded() bool         { return o.json != nil }

func (o *Note) Attachments() []*Attachment {
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

func (o *Note) InReplyTo() *Object {
	if v := o.json.Get("inReplyTo"); v != nil && v.Type() != fastjson.TypeNull {
		return &Object{v}
	}
	return nil
}
