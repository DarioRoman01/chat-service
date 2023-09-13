package main

import (
	"context"
	"fmt"
	"log"

	"github.com/DarioRoman01/delfos-chat/config"
	"github.com/DarioRoman01/delfos-chat/internal/app"
	"github.com/DarioRoman01/delfos-chat/internal/repository/mongo"
	"github.com/Netflix/go-env"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()
	var conf config.Config
	_, err := env.UnmarshalFromEnviron(&conf)
	if err != nil {
		log.Fatalf("main: env.UnmarshalFromEnviron error: %v", err)
	}

	logger := zap.Must(zap.NewProduction())
	if conf.Environment == "local" {
		logger = zap.Must(zap.NewDevelopment())
	}

	defer logger.Sync()

	mongo, err := mongo.New(&conf.Mongo)
	if err != nil {
		logger.Fatal(fmt.Sprintf("main: mongo.New error: %v", err))
	}

	defer func() {
		if err := mongo.Client().Disconnect(ctx); err != nil {
			logger.Fatal(fmt.Sprintf("main: mongo.Client.Disconnect error: %v", err))
		}
	}()

	app := app.New(&app.Config{
		Name:               "delfos-char-service",
		Ctx:                ctx,
		Logger:             logger,
		Mongo:              mongo,
		HttpPort:           conf.HttpPort,
		MessagesCollection: conf.MongoMessageCollection,
		ChannelsCollection: conf.MongoChannelCollection,
	})

	if err = app.Start(); err != nil {
		logger.Fatal(fmt.Sprintf("main: app.Start error: %v", err))
	}
}
