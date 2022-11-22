package activitypub

func (o *Outbox) TotalItems() uint { return o.json.GetUint("totalItems") }
func (o *Outbox) First() string    { return string(o.json.GetStringBytes("first")) }
func (o *Outbox) Last() string     { return string(o.json.GetStringBytes("last")) }
func (o *Outbox) Next() string     { return string(o.json.GetStringBytes("next")) }
func (o *Outbox) Prev() string     { return string(o.json.GetStringBytes("prev")) }

func (o *Outbox) Items() []item {
	items := make([]item, 0)
	for _, v := range o.json.GetArray("orderedItems") {
		items = append(items, item{v, v.GetStringBytes("id")})
	}
	return items
}
