package chat

import (
	"time"

	"github.com/DarioRoman01/chat-service/entities"
	"github.com/gofiber/contrib/websocket"
)

type client struct {
	conn *websocket.Conn
}

type channel struct {
	clients   map[*client]bool
	broadcast chan *entities.Message
	joining   chan *client
	leaving   chan *client
	close     chan struct{}
}

func NewChannel() *channel {
	return &channel{
		clients:   make(map[*client]bool),
		broadcast: make(chan *entities.Message),
		joining:   make(chan *client),
		leaving:   make(chan *client),
		close:     make(chan struct{}),
	}
}

func (c *channel) Run() {
	for {
		select {
		case client := <-c.joining:
			c.clients[client] = true

		case client := <-c.leaving:
			delete(c.clients, client)
			client.conn.Close()

		case message := <-c.broadcast:
			for client := range c.clients {
				err := client.conn.WriteJSON(message)
				if err != nil {
					return
				}
			}

		case <-c.close:
			close(c.broadcast)
			close(c.joining)
			close(c.leaving)
			close(c.close)
			return
		}
	}
}

func (c *channel) Close() {
	for client := range c.clients {
		client.conn.WriteControl(websocket.CloseMessage, []byte("closing connection"), time.Now().Add(time.Second*2))
		delete(c.clients, client)
	}

	c.close <- struct{}{}
}
