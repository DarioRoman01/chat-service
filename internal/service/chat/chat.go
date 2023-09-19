package chat

import (
	"encoding/json"
	"io"
	"time"

	"github.com/DarioRoman01/chat-service/entities"
	"github.com/DarioRoman01/chat-service/internal/repository"
	"github.com/DarioRoman01/chat-service/pkg/errors"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/oklog/ulid/v2"
)

type ChatService struct {
	repository *repository.Repository
	channels   map[string]*channel
}

func New(repository *repository.Repository) *ChatService {
	return &ChatService{
		repository: repository,
		channels:   make(map[string]*channel),
	}
}

func (c *ChatService) HandleConnection(conn *websocket.Conn, channelId string) {
	channel, err := c.getChannel(channelId)
	if err != nil {
		c.closeConnWIthErr(conn, errors.Wrap(err, "messages: ChatService.HandleConnection c.getChannel error"))
		return
	}

	client := &client{conn}
	channel.joining <- client

	for {
		message := new(entities.Message)
		connIsFinished := c.readJSON(conn, message)
		if connIsFinished {
			if _, ok := c.channels[channelId]; ok {
				// if the channel stills exists means that the client end up the connection
				channel.leaving <- client
			}
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
					"errors": errors.Wrap(err, "messages: ChatService.HandleConnection m.repository.Messages.Create error").Error(),
				})
				continue
			}

			channel.broadcast <- message

		case entities.UpdateMessage:
			message, err := c.repository.Messages.Update(message)
			if err != nil {
				conn.WriteJSON(fiber.Map{
					"errors": errors.Wrap(err, "message: ChatService.HandleConnection m.repository.Messages.Update error").Error(),
				})

				continue
			}

			channel.broadcast <- message
		}

	}
}

func (c *ChatService) CloseChannel(channelId string) error {
	channel, ok := c.channels[channelId]
	if !ok {
		return errors.Newf("messages: ChatService.Closechannel error: no running channel with id: %s", channelId)
	}

	channel.Close()
	delete(c.channels, channelId)
	return nil
}

func (c *ChatService) getChannel(channelId string) (*channel, error) {
	channel, ok := c.channels[channelId]
	if !ok {
		ok, err := c.repository.Channels.Exists(channelId)
		if err != nil {
			return nil, errors.Wrap(err, "channels: ChatService.Get c.repository.Channels.Get error")
		}

		if !ok {
			return nil, errors.Newf("channels: ChatService.Get c.repository.Channels.Exists error: channel does not exists")
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

// reads the next json mesasge and check if the connection was finished
func (c *ChatService) readJSON(conn *websocket.Conn, v interface{}) (connIsFinished bool) {
	messageType, r, err := conn.NextReader()
	if messageType == -1 { // no frames
		return true
	}

	if err != nil {
		conn.WriteJSON(fiber.Map{
			"errors": errors.Wrap(err, "messages: ChatService.readJson conn.NextReader error").Error(),
		})
	}

	err = json.NewDecoder(r).Decode(v)
	if err == io.EOF {
		// One value is expected in the message.
		err = io.ErrUnexpectedEOF
	}

	if err != nil {
		conn.WriteJSON(fiber.Map{
			"errors": errors.Wrap(err, "messages: ChatService.readJson json.NewDecoder.Decode error").Error(),
		})
	}

	return false
}
