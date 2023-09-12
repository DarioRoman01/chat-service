package messages

import (
	"github.com/DarioRoman01/delfos-chat/entities"
	"github.com/gofiber/contrib/websocket"
)

type Client struct {
	conn *websocket.Conn
}

type ChatRoom struct {
	clients   map[*Client]bool
	broadcast chan *entities.Message
	joining   chan *Client
	leaving   chan *Client
}

func NewChatRoom() *ChatRoom {
	return &ChatRoom{
		clients:   make(map[*Client]bool),
		broadcast: make(chan *entities.Message),
		joining:   make(chan *Client),
		leaving:   make(chan *Client),
	}
}

func (c *ChatRoom) Run() {
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
		}
	}
}

func (c *ChatRoom) Close() {
	for client := range c.clients {
		client.conn.Close()
		delete(c.clients, client)
	}
}
