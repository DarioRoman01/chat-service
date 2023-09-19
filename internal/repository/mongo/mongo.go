package mongo

import (
	"context"
	"fmt"

	"github.com/DarioRoman01/chat-service/config"
	"github.com/DarioRoman01/chat-service/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func New(conf *config.Mongo) (*mongo.Database, error) {
	mongoURI := fmt.Sprintf(
		"mongodb://%s:%s@%s/%s?directConnection=true&serverSelectionTimeoutMS=2000&authSource=admin",
		conf.MongoUsername,
		conf.MongoPassword,
		conf.MongoHost,
		conf.MongoDBName,
	)

	client, err := mongo.Connect(
		context.Background(),
		options.Client().ApplyURI(mongoURI),
	)
	if err != nil {
		return nil, errors.Wrap(err, "mongo: GetDatabaseConnection mongo.Connect error")
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "mongo: GetDatabaseConnection client.Ping error")
	}

	return client.Database(conf.MongoDBName), nil
}
