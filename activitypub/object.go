package activitypub

import "github.com/valyala/fastjson"

func (o *Object) Type() string { return string(o.json.GetStringBytes("type")) }

func (o *Object) ToNote() *Note             { return &Note{o.json, o.ID} }
func (o *Object) ToPerson() *Person         { return &Person{o.json, o.ID} }
func (o *Object) ToCollection() *Collection { return &Collection{o.json, o.ID} }

func NewObject(v *fastjson.Value) *Object {
	return &Object{v, v.GetStringBytes("id")}
}
