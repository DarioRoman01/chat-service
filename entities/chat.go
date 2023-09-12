package entities

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type MessageType int

const (
	SendMessage   MessageType = iota
	UpdateMessage MessageType = 1
)

type Message struct {
	Id ulid.ULID `json:"id" bson:"id"` // id of the message

	Type MessageType `json:"-" bson:"-"` // type of the message that is sent

	Content string `json:"content" bson:"content"` // text content inside the message

	Media string `json:"media" bson:"media"` // possible media inside the message, not defined yet how to handle this

	SentBy string `json:"sentBy" bson:"sentBy"` // who sent the message

	RecievedBy []string `json:"recievedBy" bson:"recievedBy"` // who recieved the message

	Seenby []string `json:"seenBy" bson:"seenBy"`

	Channel string `json:"channel" bson:"channel"` // channel id where the message was sent

	Tournament string `json:"tournament" bson:"tournament"` // possible tournament where the message comes from

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"` // creation timestamp
}

// TODO: maybe a possible chat struct?
