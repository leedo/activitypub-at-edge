package activitypub

func (o *Object) Type() string { return string(o.json.GetStringBytes("type")) }

func (o *Object) ToNote() *Note             { return &Note{o.json, o.ID} }
func (o *Object) ToPerson() *Person         { return &Person{o.json, o.ID} }
func (o *Object) ToCollection() *Collection { return &Collection{o.json, o.ID} }

// lazy object if json is a string and not an object?
