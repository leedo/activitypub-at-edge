package activitypub

import "github.com/valyala/fastjson"

type activityPubObject struct {
	json *fastjson.Value
	ID   []byte
}

type Person activityPubObject
type Outbox activityPubObject
type Item activityPubObject
type Object activityPubObject

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
