package channels

import (
	"github.com/DarioRoman01/chat-service/entities"
	"github.com/DarioRoman01/chat-service/internal/repository"
	"github.com/DarioRoman01/chat-service/pkg/errors"
)

type ChannelService struct {
	repository *repository.Repository
}

func New(repository *repository.Repository) *ChannelService {
	return &ChannelService{
		repository: repository,
	}
}

func (c *ChannelService) Create(channelId string) error {
	_, err := c.repository.Channels.Create(&entities.Channel{
		Id:         channelId,
		Tournament: "someTournament",
		Public:     false,
	})

	if err != nil {
		return errors.Wrap(err, "channels: ChannelService.Create c.repository.Channels.Create error")
	}

	return nil
}
