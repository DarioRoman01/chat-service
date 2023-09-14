package service

import (
	"github.com/DarioRoman01/delfos-chat/entities"
	"github.com/DarioRoman01/delfos-chat/internal/repository"
	"github.com/DarioRoman01/delfos-chat/internal/service/channels"
	"github.com/DarioRoman01/delfos-chat/internal/service/chat"
	"github.com/DarioRoman01/delfos-chat/internal/service/messages"
	"github.com/gofiber/contrib/websocket"
)

type MessageServices interface {
	Get(string) ([]*entities.Message, error)
}

type ChannelService interface {
	Create(string) error
}

type ChatService interface {
	HandleConnection(*websocket.Conn, string)
	CloseChannel(string) error
}

type Service struct {
	Messages MessageServices
	Channels ChannelService
	Chats    ChatService
}

func New(repository *repository.Repository) *Service {
	return &Service{
		Messages: messages.New(repository),
		Channels: channels.New(repository),
		Chats:    chat.New(repository),
	}
}
