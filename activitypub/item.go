package activitypub

func (o *Item) Actor() string { return string(o.json.GetStringBytes("actor")) }
