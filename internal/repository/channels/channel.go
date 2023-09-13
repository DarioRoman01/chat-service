package channels

import (
	"context"

	"github.com/DarioRoman01/delfos-chat/entities"
	"github.com/DarioRoman01/delfos-chat/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ChannelRepository struct {
	mongo      *mongo.Database
	collection *mongo.Collection
	ctx        context.Context
}

type Config struct {
	Mongo      *mongo.Database
	Collection string
	Ctx        context.Context
}

func New(conf *Config) *ChannelRepository {
	return &ChannelRepository{
		mongo:      conf.Mongo,
		collection: conf.Mongo.Collection(conf.Collection),
		ctx:        conf.Ctx,
	}
}

func (c *ChannelRepository) Create(channel *entities.Channel) (*entities.Channel, error) {
	_, err := c.collection.InsertOne(c.ctx, channel)
	if err != nil {
		return nil, errors.Wrap(err, "channel: Create c.collection.InsertOne error")
	}

	return channel, nil
}

func (c *ChannelRepository) Exists(channelId string) (bool, error) {
	result := c.collection.FindOne(c.ctx, bson.M{"id": channelId})

	var channel entities.Channel
	if err := result.Decode(&channel); err != nil {
		return false, errors.Wrap(err, "channel: Get result.Decode error")
	}

	return channel.Id == channelId, nil
}
