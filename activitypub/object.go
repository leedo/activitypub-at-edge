package activitypub

func (o *Object) Type() string         { return string(o.json.GetStringBytes("type")) }
func (o *Object) InReplyTo() string    { return string(o.json.GetStringBytes("inReplyTo")) }
func (o *Object) URL() string          { return string(o.json.GetStringBytes("url")) }
func (o *Object) Content() []byte      { return o.json.GetStringBytes("content") }
func (o *Object) Published() string    { return string(o.json.GetStringBytes("published")) }
func (o *Object) AttributedTo() string { return string(o.json.GetStringBytes("attributedTo")) }

func (i *item) Object() *object {
	o := &Object{}
	if i.json.GetObject("object") != nil {
		o.json = i.json.Get("object")
		o.id = o.json.GetStringBytes("id")
	} else {
		o.id = i.json.GetStringBytes("object")
	}
	return o
}
