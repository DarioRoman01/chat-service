package chat

import (
	"encoding/json"
	"io"
	"time"

	"github.com/DarioRoman01/delfos-chat/entities"
	"github.com/DarioRoman01/delfos-chat/internal/repository"
	"github.com/DarioRoman01/delfos-chat/pkg/errors"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/oklog/ulid/v2"
)

type ChatService struct {
	repository *repository.Repository
	channels   map[string]*Channel
}

func New(repository *repository.Repository) *ChatService {
	return &ChatService{
		repository: repository,
		channels:   make(map[string]*Channel),
	}
}

func (c *ChatService) HandleConnection(conn *websocket.Conn, channelId string) {
	channel, err := c.getChannel(channelId)
	if err != nil {
		c.closeConnWIthErr(conn, errors.Wrap(err, "messages: HandleConnection error"))
		return
	}

	client := &Client{Conn: conn}
	channel.Joining <- client
	defer func() { channel.Leaving <- client }()

	for {
		message := new(entities.Message)
		connIsFinished := c.readJSON(conn, message)
		if connIsFinished {
			return
		}

		switch message.Type {
		case entities.SendMessage:
			message, err = c.repository.Messages.Create(&entities.Message{
				Id:         ulid.Make(),
				Content:    message.Content,
				SentBy:     message.SentBy,
				RecievedBy: make([]string, 0),
				Seenby:     make([]string, 0),
				Media:      message.Media,
				Channel:    channelId,
				CreatedAt:  time.Now(),
			})

			if err != nil {
				conn.WriteJSON(fiber.Map{
					"errors": errors.Wrap(err, "messages: HandleConnection m.repository.Messages.Create error").Error(),
				})
				continue
			}

			channel.Broadcast <- message

		case entities.UpdateMessage:
			message, err := c.repository.Messages.Update(message)
			if err != nil {
				conn.WriteJSON(fiber.Map{
					"errors": errors.Wrap(err, "message: HandleConnection m.repository.Messages.Update error").Error(),
				})

				continue
			}

			channel.Broadcast <- message
		}

	}
}

func (c *ChatService) CloseChannel(channelId string) error {
	channel, ok := c.channels[channelId]
	if !ok {
		return errors.Newf("messages: Closechannel error: no running channel with id: %s", channelId)
	}

	channel.Close()
	delete(c.channels, channelId)
	return nil
}

func (c *ChatService) getChannel(channelId string) (*Channel, error) {
	channel, ok := c.channels[channelId]
	if !ok {
		ok, err := c.repository.Channels.Exists(channelId)
		if err != nil {
			return nil, errors.Wrap(err, "channels: Get c.repository.Channels.Get error")
		}

		if !ok {
			return nil, errors.Newf("channels: Get c.repository.Channels.Exists error: channel does not exists")
		}

		channel = NewChannel()
		c.channels[channelId] = channel
		go channel.Run()
	}

	return channel, nil
}

func (c *ChatService) closeConnWIthErr(conn *websocket.Conn, err error) {
	conn.WriteJSON(fiber.Map{"errors": err.Error()})
	conn.Close()
}

// reads the next json mesasge and check if the connection was finished by the client
func (c *ChatService) readJSON(conn *websocket.Conn, v interface{}) (connIsFinished bool) {
	messageType, r, err := conn.NextReader()
	if messageType == -1 { // no frames
		return true
	}

	if err != nil {
		conn.WriteJSON(fiber.Map{
			"errors": errors.Wrap(err, "messages: HandleConnection conn.ReadJson error").Error(),
		})
	}

	err = json.NewDecoder(r).Decode(v)
	if err == io.EOF {
		// One value is expected in the message.
		err = io.ErrUnexpectedEOF
	}

	if err != nil {
		conn.WriteJSON(fiber.Map{
			"errors": errors.Wrap(err, "messages: HandleConnection conn.ReadJson error").Error(),
		})
	}

	return false
}
