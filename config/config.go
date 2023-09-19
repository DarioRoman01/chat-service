package config

import (
	"github.com/DarioRoman01/chat-service/pkg/errors"
	"github.com/Netflix/go-env"
)

type (
	Config struct {
		App
		Http
		Mongo
	}

	App struct {
		Environment      string `env:"CHAT_SERVICE_ENVIRONMENT,required=true"`
		MaxParallelChats string `env:"MAX_PARALLEL_CHATS,required=true"`
	}

	Http struct {
		HttpPort string `env:"CHAT_SERICE_PORT,required=true"`
	}

	Mongo struct {
		MongoUsername          string `env:"MONGODB_USERNAME,required=true"`
		MongoPassword          string `env:"MONGODB_PASSWORD,required=true"`
		MongoHost              string `env:"MONGODB_HOST,required=true"`
		MongoDBName            string `env:"MONGODB_DATABASE,required=true"`
		MongoMessageCollection string `env:"MONGODB_MESSAGES_COLLECTION,required=true"`
		MongoChannelCollection string `env:"MONGODB_CHANNELS_COLLECTION,required=true"`
	}
)

func New() (*Config, error) {
	var conf Config
	_, err := env.UnmarshalFromEnviron(&conf)
	if err != nil {
		return nil, errors.Wrap(err, "config: New env.UnmarshalFromEnviron error")
	}

	return &conf, nil
}
