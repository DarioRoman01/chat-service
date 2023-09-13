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
	Id         ulid.ULID   `json:"id" bson:"id"`
	Type       MessageType `json:"type" bson:"-"`
	Content    string      `json:"content" bson:"content"`
	Media      string      `json:"media" bson:"media"`
	SentBy     string      `json:"sentBy" bson:"sentBy"`
	RecievedBy []string    `json:"recievedBy" bson:"recievedBy"`
	Seenby     []string    `json:"seenBy" bson:"seenBy"`
	Channel    string      `json:"channel" bson:"channel"`
	CreatedAt  time.Time   `json:"createdAt" bson:"createdAt"`
}

type Channel struct {
	Id         string `json:"id" bson:"id"`
	Tournament string `json:"tournament" bson:"tournament"`
	Public     bool   `json:"public" bson:"public"`
}
