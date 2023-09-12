package repository

import (
	"context"

	"github.com/DarioRoman01/delfos-chat/internal/repository/messages"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	Messages MessageRepository
}

type Config struct {
	Mongo              *mongo.Database
	Ctx                context.Context
	MessagesCollection string
}

func New(conf *Config) *Repository {
	return &Repository{
		Messages: messages.New(&messages.Config{
			Mongo:      conf.Mongo,
			Collection: conf.MessagesCollection,
			Ctx:        conf.Ctx,
		}),
	}
}
