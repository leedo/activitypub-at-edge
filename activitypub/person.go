package activitypub

func (p *Person) Outbox() string { return string(p.json.GetStringBytes("outbox")) }
func (p *Person) Name() string   { return string(p.json.GetStringBytes("name")) }

func (p *Person) Icon() image {
	return image{string(p.json.Get("icon").GetStringBytes("url"))}
}
