package messages

import (
	"github.com/DarioRoman01/delfos-chat/entities"
	"github.com/DarioRoman01/delfos-chat/internal/repository"
	"github.com/DarioRoman01/delfos-chat/pkg/errors"
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
		return nil, errors.Wrap(err, "messages: Get m.repository.Messages.Get error")
	}

	return messages, nil
}
