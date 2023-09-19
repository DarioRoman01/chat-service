package repository

import "github.com/DarioRoman01/chat-service/entities"

type MessageRepository interface {
	// Creates a new message record in the database
	Create(*entities.Message) (*entities.Message, error)

	// Get the latest messages from a channel
	Get(string) ([]*entities.Message, error)

	// Updates the given record message
	Update(*entities.Message) (*entities.Message, error)
}

type ChannelRepository interface {
	Create(*entities.Channel) (*entities.Channel, error)
	Exists(string) (bool, error)
}
