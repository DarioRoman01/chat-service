package repository

import (
	"context"

	"github.com/DarioRoman01/chat-service/internal/repository/channels"
	"github.com/DarioRoman01/chat-service/internal/repository/messages"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	Messages MessageRepository
	Channels ChannelRepository
}

type Config struct {
	Mongo              *mongo.Database
	Ctx                context.Context
	MessagesCollection string
	ChannelCollection  string
}

func New(conf *Config) *Repository {
	return &Repository{
		Messages: messages.New(&messages.Config{
			Mongo:      conf.Mongo,
			Collection: conf.MessagesCollection,
			Ctx:        conf.Ctx,
		}),
		Channels: channels.New(&channels.Config{
			Mongo:      conf.Mongo,
			Collection: conf.ChannelCollection,
			Ctx:        conf.Ctx,
		}),
	}
}
