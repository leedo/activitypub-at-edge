package activitypub

import "github.com/valyala/fastjson"

type Object struct {
	json *fastjson.Value
	ID   []byte
}

type Person Object
type Collection Object
type CollectionItem Object
type Note Object

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
