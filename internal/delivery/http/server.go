package http

import (
	"github.com/DarioRoman01/delfos-chat/internal/delivery/http/routes"
	"github.com/DarioRoman01/delfos-chat/internal/service"
	"github.com/DarioRoman01/delfos-chat/pkg/http"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.uber.org/zap"
)

func New(service *service.Service, logger *zap.Logger) *fiber.App {
	server := fiber.New()
	server.Use(cors.New())
	server.Use(http.LogMiddleware(logger))
	router := server.Group("delfos/api/v1")
	routes.New(service, router)
	return server
}
