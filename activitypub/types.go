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
