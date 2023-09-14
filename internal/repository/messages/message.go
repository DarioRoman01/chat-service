package messages

import (
	"context"

	"github.com/DarioRoman01/delfos-chat/entities"
	"github.com/DarioRoman01/delfos-chat/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MessageRepository struct {
	mongo      *mongo.Database
	collection *mongo.Collection
	ctx        context.Context
}

type Config struct {
	Mongo      *mongo.Database
	Collection string
	Ctx        context.Context
}

func New(conf *Config) *MessageRepository {
	return &MessageRepository{
		mongo:      conf.Mongo,
		collection: conf.Mongo.Collection(conf.Collection),
		ctx:        conf.Ctx,
	}
}

func (m *MessageRepository) Create(message *entities.Message) (*entities.Message, error) {
	_, err := m.collection.InsertOne(m.ctx, message, nil)
	if err != nil {
		return nil, errors.Wrap(err, "messages: MessageRepository.Create m.collection.InsertOne error")
	}

	return message, nil
}

func (m *MessageRepository) Get(channel string) ([]*entities.Message, error) {
	cursor, err := m.collection.Find(m.ctx, bson.M{"channel": channel})
	if err != nil {
		return nil, errors.Wrap(err, "messages: MessageRepository.Get m.collection.Find error")
	}

	defer cursor.Close(m.ctx)
	messages := make([]*entities.Message, 0, cursor.RemainingBatchLength())
	if err := cursor.All(m.ctx, &messages); err != nil {
		return nil, errors.Wrap(err, "messages: MessageRepository.Get cursor.All error")
	}

	return messages, nil
}

func (m *MessageRepository) Update(message *entities.Message) (*entities.Message, error) {
	updateResult, err := m.collection.UpdateOne(m.ctx, bson.M{"id": message.Id.String()}, message)
	if err != nil {
		return nil, errors.Wrap(err, "messages: Update m.collection.UpdateOne error")
	}

	if updateResult.ModifiedCount < 1 {
		return nil, errors.New("messages: MessageRepository.Update updateResult.ModifiedCount is less than one: no record found")
	}

	res := m.collection.FindOne(m.ctx, bson.M{"id": message.Id})
	updatedMessage := new(entities.Message)
	if err := res.Decode(updatedMessage); err != nil {
		return nil, errors.Wrap(err, "messages: MessageRepository.Update res.Decode error")
	}

	return updatedMessage, nil
}
