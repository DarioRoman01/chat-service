package messages

import (
	"github.com/DarioRoman01/chat-service/entities"
	"github.com/DarioRoman01/chat-service/internal/repository"
	"github.com/DarioRoman01/chat-service/pkg/errors"
)

type MessageService struct {
	repository *repository.Repository
}

func New(repo *repository.Repository) *MessageService {
	return &MessageService{
		repository: repo,
	}
}

func (m *MessageService) Get(channelId string) ([]*entities.Message, error) {
	messages, err := m.repository.Messages.Get(channelId)
	if err != nil {
		return nil, errors.Wrap(err, "messages: MessageService.Get m.repository.Messages.Get error")
	}

	return messages, nil
}
