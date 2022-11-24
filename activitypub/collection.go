package activitypub

func (c *Collection) ID() string                 { return string(c.json.GetStringBytes("id")) }
func (c *Collection) TotalCollectionItems() uint { return c.json.GetUint("totalCollectionItems") }
func (c *Collection) First() string              { return string(c.json.GetStringBytes("first")) }
func (c *Collection) Last() string               { return string(c.json.GetStringBytes("last")) }
func (c *Collection) Next() string               { return string(c.json.GetStringBytes("next")) }
func (c *Collection) Prev() string               { return string(c.json.GetStringBytes("prev")) }

func (c *Collection) CollectionItems() []*Object {
	items := make([]*Object, 0)
	for _, v := range c.json.GetArray("orderedItems") {
		items = append(items, &Object{v})
	}
	return items
}
