package routes

import (
	"github.com/DarioRoman01/delfos-chat/internal/service"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ChatDelivery struct {
	service *service.Service
}

func NewChatDelivery(service *service.Service, router fiber.Router) {
	chatDelivery := &ChatDelivery{
		service: service,
	}

	route := router.Group("/chat")
	chatDelivery.Create(route)
	chatDelivery.Join(route)
	chatDelivery.Close(route)
	chatDelivery.GetChannelMessages(route)
}

func (ch *ChatDelivery) Create(router fiber.Router) {
	router.Post("/create", func(c *fiber.Ctx) error {
		id := uuid.NewString()
		if err := ch.service.Channels.Create(id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"errors": err.Error(),
			})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"id": id,
		})
	})
}

func (ch *ChatDelivery) Join(router fiber.Router) {
	router.Get("/:id/join", websocket.New(func(conn *websocket.Conn) {
		ch.service.Chats.HandleConnection(conn, conn.Params("id"))
	}))
}

func (ch *ChatDelivery) GetChannelMessages(router fiber.Router) {
	router.Get("/:id/messages", func(c *fiber.Ctx) error {
		messages, err := ch.service.Messages.Get(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"errors": err,
			})
		}

		return c.JSON(messages)
	})
}

func (ch *ChatDelivery) Close(router fiber.Router) {
	router.Delete("/:id/close", func(c *fiber.Ctx) error {
		err := ch.service.Chats.CloseChannel(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"errors": err,
			})
		}

		return c.SendStatus(fiber.StatusNoContent)
	})
}
