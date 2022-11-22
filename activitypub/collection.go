package activitypub

import "github.com/valyala/fastjson"

func (o *Collection) TotalCollectionItems() uint { return o.json.GetUint("totalCollectionItems") }
func (o *Collection) First() string              { return string(o.json.GetStringBytes("first")) }
func (o *Collection) Last() string               { return string(o.json.GetStringBytes("last")) }
func (o *Collection) Next() string               { return string(o.json.GetStringBytes("next")) }
func (o *Collection) Prev() string               { return string(o.json.GetStringBytes("prev")) }

func (o *Collection) CollectionItems() []*CollectionItem {
	items := make([]*CollectionItem, 0)
	for _, v := range o.json.GetArray("orderedItems") {
		items = append(items, &CollectionItem{v, v.GetStringBytes("id")})
	}
	return items
}

func NewCollection(v *fastjson.Value) *Collection {
	return &Collection{v, v.GetStringBytes("id")}
}
