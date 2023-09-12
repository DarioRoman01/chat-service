package routes

import (
	"github.com/DarioRoman01/delfos-chat/internal/service"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type MessagesDelivery struct {
	service *service.Service
}

func New(service *service.Service, router fiber.Router) {
	delivery := &MessagesDelivery{
		service: service,
	}

	route := router.Group("/chat")
	delivery.CreateChatRoom(route)
	delivery.JoinChatRoom(route)
	delivery.CloseChatRoom(route)
}

func (m *MessagesDelivery) CreateChatRoom(router fiber.Router) {
	router.Post("/create", func(c *fiber.Ctx) error {
		id := uuid.NewString()
		m.service.Messages.CreateChannel(id)
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"id": id,
		})
	})
}

func (m *MessagesDelivery) JoinChatRoom(router fiber.Router) {
	router.Get("/:id/ws", websocket.New(func(c *websocket.Conn) {
		m.service.Messages.HandleConnection(c, c.Params("id"))
	}))
}

func (m *MessagesDelivery) CloseChatRoom(router fiber.Router) {
	router.Delete("/:id/delete", func(c *fiber.Ctx) error {
		err := m.service.Messages.CloseChatRoom(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"errors": err,
			})
		}

		return c.SendStatus(fiber.StatusNoContent)
	})
}
