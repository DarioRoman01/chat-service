package http

import (
	"github.com/DarioRoman01/chat-service/internal/delivery/http/routes"
	"github.com/DarioRoman01/chat-service/internal/service"
	"github.com/DarioRoman01/chat-service/pkg/http/middlewares"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.uber.org/zap"
)

func New(service *service.Service, logger *zap.Logger) *fiber.App {
	server := fiber.New()
	server.Use(cors.New())
	server.Use(middlewares.LogMiddleware(logger))

	router := server.Group("delfos/api/v1")
	routes.NewChatDelivery(service, router)
	return server
}
