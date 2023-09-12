package service

import (
	"github.com/DarioRoman01/delfos-chat/internal/repository"
	"github.com/DarioRoman01/delfos-chat/internal/service/messages"
	"github.com/gofiber/contrib/websocket"
)

type MessageServices interface {
	HandleConnection(*websocket.Conn, string)
	CreateChannel(string)
	CloseChatRoom(string) error
}

type Service struct {
	Messages MessageServices
}

func New(repository *repository.Repository) *Service {
	return &Service{
		Messages: messages.New(repository),
	}
}
