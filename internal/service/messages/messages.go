package messages

import (
	"time"

	"github.com/DarioRoman01/delfos-chat/entities"
	"github.com/DarioRoman01/delfos-chat/internal/repository"
	"github.com/DarioRoman01/delfos-chat/pkg/errors"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/oklog/ulid/v2"
)

type MessageService struct {
	repository *repository.Repository
	chatRooms  map[string]*ChatRoom // TODO: how to handle TTL for chat rooms
}

func New(repo *repository.Repository) *MessageService {
	return &MessageService{
		repository: repo,
		chatRooms:  make(map[string]*ChatRoom),
	}
}

func (m *MessageService) HandleConnection(conn *websocket.Conn, channel string) {
	chatRoom, ok := m.chatRooms[channel]
	if !ok {
		m.closeConnWIthErr(conn, errors.New("chatRoom with id does not exists"))
		return
	}

	messages, err := m.repository.Messages.Get(channel)
	if err != nil {
		m.closeConnWIthErr(conn, errors.Wrap(err, "messages: HandleConnection m.repository.Messages.Get error"))
		return
	}

	if len(messages) != 0 {
		for message := range messages {
			if err := conn.WriteJSON(message); err != nil {
				m.closeConnWIthErr(conn, errors.Wrap(err, "messages: HandleConnection conn.WriteJSON error"))
				return
			}
		}
	}

	client := &Client{conn}
	chatRoom.joining <- client
	defer func() { chatRoom.leaving <- client }()

	for {
		message := new(entities.Message)
		if err := conn.ReadJSON(message); err != nil {
			conn.WriteJSON(fiber.Map{
				"errors": errors.Wrap(err, "messages: HandleConnection conn.ReadJson error").Error(),
			})
			continue
		}

		switch message.Type {
		case entities.SendMessage:
			message, err = m.repository.Messages.Create(&entities.Message{
				Id:         ulid.Make(),
				Content:    message.Content,
				SentBy:     message.SentBy,
				RecievedBy: make([]string, 0),
				Seenby:     make([]string, 0),
				Media:      message.Media,
				Channel:    channel,
				Tournament: message.Tournament,
				CreatedAt:  time.Now(),
			})

			if err != nil {
				conn.WriteJSON(fiber.Map{
					"errors": errors.Wrap(err, "messages: HandleConnection m.repository.Messages.Create error").Error(),
				})
				continue
			}

			chatRoom.broadcast <- message

		case entities.UpdateMessage:
			message, err := m.repository.Messages.Update(message)
			if err != nil {
				conn.WriteJSON(fiber.Map{
					"errors": errors.Wrap(err, "message: HandleConnection m.repository.Messages.Update error").Error(),
				})

				continue
			}

			chatRoom.broadcast <- message
		}

	}
}

func (m *MessageService) closeConnWIthErr(conn *websocket.Conn, err error) {
	conn.WriteJSON(fiber.Map{"errors": err.Error()})
	conn.Close()
}

func (m *MessageService) CreateChannel(channelId string) {
	chatRoom := NewChatRoom()
	m.chatRooms[channelId] = chatRoom
	go chatRoom.Run()
}

func (m *MessageService) CloseChatRoom(channelId string) error {
	chatRoom, ok := m.chatRooms[channelId]
	if !ok {
		return errors.New("messages: CloseChatRoom error: chatRoom with id does not exists")
	}

	chatRoom.Close()
	delete(m.chatRooms, channelId)
	return nil
}
