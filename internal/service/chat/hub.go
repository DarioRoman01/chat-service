package chat

import (
	"time"

	"github.com/DarioRoman01/delfos-chat/entities"
	"github.com/gofiber/contrib/websocket"
)

type Client struct {
	Conn *websocket.Conn
}

type Channel struct {
	clients   map[*Client]bool
	Broadcast chan *entities.Message
	Joining   chan *Client
	Leaving   chan *Client
}

func NewChannel() *Channel {
	return &Channel{
		clients:   make(map[*Client]bool),
		Broadcast: make(chan *entities.Message),
		Joining:   make(chan *Client),
		Leaving:   make(chan *Client),
	}
}

func (c *Channel) Run() {
	for {
		select {
		case client := <-c.Joining:
			c.clients[client] = true

		case client := <-c.Leaving:
			delete(c.clients, client)
			client.Conn.Close()

		case message := <-c.Broadcast:
			for client := range c.clients {
				err := client.Conn.WriteJSON(message)
				if err != nil {
					return
				}
			}
		}
	}
}

func (c *Channel) Close() {
	for client := range c.clients {
		client.Conn.WriteControl(websocket.CloseMessage, []byte("closing connection"), time.Now().Add(time.Second*2))
		delete(c.clients, client)
	}
}
