package activitypub

import "github.com/valyala/fastjson"

func (o *Object) Type() string { return string(o.json.GetStringBytes("type")) }

func (o *Object) ToNote() *Note             { return &Note{o.json} }
func (o *Object) ToPerson() *Person         { return &Person{o.json} }
func (o *Object) ToCollection() *Collection { return &Collection{o.json} }
func (o *Object) ToActivity() *Activity     { return &Activity{o.json} }
func (o *Object) ToQuestion() *Question     { return &Question{o.json} }

func (o *Object) IsLoaded() bool {
	return o.json.Type() == fastjson.TypeObject
}

func (o *Object) ID() string {
	if o.IsLoaded() {
		return string(o.json.GetStringBytes("id"))
	}
	return string(o.json.GetStringBytes())
}
