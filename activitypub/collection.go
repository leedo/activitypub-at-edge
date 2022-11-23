package activitypub

import "github.com/valyala/fastjson"

func (c *Collection) TotalCollectionItems() uint { return o.json.GetUint("totalCollectionItems") }
func (c *Collection) First() string              { return string(o.json.GetStringBytes("first")) }
func (c *Collection) Last() string               { return string(o.json.GetStringBytes("last")) }
func (c *Collection) Next() string               { return string(o.json.GetStringBytes("next")) }
func (c *Collection) Prev() string               { return string(o.json.GetStringBytes("prev")) }

func (c *Collection) CollectionItems() []*CollectionItem {
	items := make([]*CollectionItem, 0)
	for _, v := range c.json.GetArray("orderedItems") {
		items = append(items, &CollectionItem{v, v.GetStringBytes("id")})
	}
	return items
}

func NewCollection(v *fastjson.Value) *Collection {
	return &Collection{v, v.GetStringBytes("id")}
}
