package activitypub

import "github.com/valyala/fastjson"

const (
	PersonType                = "Person"
	NoteType                  = "Note"
	OrderedCollectionPageType = "OrderedCollectionPage"
	OrderedCollectionType     = "OrderedCollection"
	CreateType                = "Create"
	AnnounceType              = "Announce"
	QuestionType              = "Question"
)

type Object struct {
	json *fastjson.Value
}

type Person Object
type Collection Object
type Note Object
type Activity Object
type Question Object

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
