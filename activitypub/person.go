package activitypub

import "github.com/valyala/fastjson"

func (p *Person) Outbox() string { return string(p.json.GetStringBytes("outbox")) }
func (p *Person) Name() string   { return string(p.json.GetStringBytes("name")) }

func (p *Person) Icon() Image {
	return Image{string(p.json.Get("icon").GetStringBytes("url"))}
}

func NewPerson(v *fastjson.Value) *Person {
	return &Person{v, v.GetStringBytes("id")}
}
