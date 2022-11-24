package activitypub

func (c *Collection) TotalCollectionItems() uint { return c.json.GetUint("totalCollectionItems") }
func (c *Collection) First() string              { return string(c.json.GetStringBytes("first")) }
func (c *Collection) Last() string               { return string(c.json.GetStringBytes("last")) }
func (c *Collection) Next() string               { return string(c.json.GetStringBytes("next")) }
func (c *Collection) Prev() string               { return string(c.json.GetStringBytes("prev")) }

func (c *Collection) Iterator() *CollectionIterator {
	v := c.json.GetArray("orderedItems")
	return &CollectionIterator{v, 0, len(v)}
}

func (i *CollectionIterator) Next() bool {
	return i.pos < i.count
}

func (i *CollectionIterator) Activity() *Activity {
	return nil
}

func (c *Collection) CollectionItems() []*Object {
	items := make([]*Object, 0)
	for _, v := range c.json.GetArray("orderedItems") {
		items = append(items, &Object{v})
	}
	return items
}
