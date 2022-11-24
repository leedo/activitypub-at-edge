package activitypub

func (a *Activity) ID() string      { return string(a.json.GetStringBytes("id")) }
func (a *Activity) Actor() string   { return string(a.json.GetStringBytes("actor")) }
func (a *Activity) Type() string    { return string(a.json.GetStringBytes("type")) }
func (a *Activity) Object() *Object { return &Object{a.json.Get("object")} }
