package activitypub

import "github.com/valyala/fastjson"

const (
	PersonType                = "Person"
	NoteType                  = "Note"
	OrderedCollectionPageType = "OrderedCollectionPage"
	OrderedCollectionType     = "OrderedCollection"
	CreateType                = "Create"
	AnnounceType              = "Announce"
)

type Object struct {
	json *fastjson.Value
}

type CollectionIterator struct {
	v     []*fastjson.Value
	pos   int
	count int
}

type Person Object
type Collection Object
type Note Object
type Activity Object

type Image struct {
	URL string
}

type Attachment struct {
	Type      string
	MediaType string
	URL       string
	Width     uint
	Height    uint
	Name      string
}
