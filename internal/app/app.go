package app

import (
	"context"
	"fmt"

	"github.com/DarioRoman01/delfos-chat/internal/delivery/http"
	"github.com/DarioRoman01/delfos-chat/internal/repository"
	"github.com/DarioRoman01/delfos-chat/internal/service"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type Application struct {
	name string

	logger *zap.Logger

	httpServer *fiber.App

	httpPort string
}

type Config struct {
	Name string

	Ctx context.Context

	Logger *zap.Logger

	Mongo *mongo.Database

	MessagesCollection string

	ChannelsCollection string

	HttpPort string
}

func New(conf *Config) *Application {
	repository := repository.New(&repository.Config{
		Mongo:              conf.Mongo,
		Ctx:                conf.Ctx,
		MessagesCollection: conf.MessagesCollection,
		ChannelCollection:  conf.ChannelsCollection,
	})

	svc := service.New(repository)
	httpServer := http.New(svc, conf.Logger)

	return &Application{
		name:       conf.Name,
		logger:     conf.Logger,
		httpServer: httpServer,
		httpPort:   conf.HttpPort,
	}
}

func (app *Application) Start() error {
	return app.httpServer.Listen(fmt.Sprintf(":%s", app.httpPort))
}
