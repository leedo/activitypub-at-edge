package activitypub

func (p *Person) ID() string     { return string(p.json.GetStringBytes("id")) }
func (p *Person) Outbox() string { return string(p.json.GetStringBytes("outbox")) }
func (p *Person) Name() string   { return string(p.json.GetStringBytes("name")) }

func (p *Person) Icon() Image {
	return Image{string(p.json.Get("icon").GetStringBytes("url"))}
}

func (p *Person) Image() Image {
	return Image{string(p.json.Get("image").GetStringBytes("url"))}
}
