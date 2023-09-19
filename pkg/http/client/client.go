package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/DarioRoman01/chat-service/entities"
	"github.com/DarioRoman01/chat-service/pkg/errors"
	"github.com/fasthttp/websocket"
)

// API client for testing purposes like e2e or integration
type Client struct {
	baseUrl    string
	httpClient *http.Client
}

func New(baseUrl string) *Client {
	return &Client{
		baseUrl:    baseUrl,
		httpClient: http.DefaultClient,
	}
}

func (c *Client) GetMessages(chatId string) ([]*entities.Message, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s/messages", c.baseUrl, chatId), nil)
	if err != nil {
		return nil, errors.Wrap(err, "client: Client.GetMessages http.NewRequest error")
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "client: Client.GetMessages httpClient.Do error")
	}

	messages := make([]*entities.Message, 0)
	err = json.NewDecoder(res.Body).Decode(&messages)
	if err != nil {
		return nil, errors.Wrap(err, "client: Client.GetMessages json.NewDecoder.Decode error")
	}

	return messages, nil
}

func (c *Client) CreateChannel() (*entities.Channel, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/chat/create", c.baseUrl), nil)
	if err != nil {
		return nil, errors.Wrap(err, "client: Client.CreateChannel http.NewRequest error")
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "client: Client.CreateChannel httpClient.Do error")
	}

	var channel entities.Channel
	err = json.NewDecoder(res.Body).Decode(&channel)
	if err != nil {
		return nil, errors.Wrap(err, "client: Client.CreateChannel json.NewDecoder.Decode error")
	}

	return &channel, nil
}

func (c *Client) WsConn(channelId string) (chan *entities.Message, chan struct{}, error) {
	messages := make(chan *entities.Message)
	closeCh := make(chan struct{})

	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("%s/chat/%s/join", c.baseUrl, channelId), http.Header{})
	if err != nil {
		return nil, nil, errors.Wrap(err, "client: Client.WsConn websocket.DefaultDialer.Dial error")
	}

	go func() {
		for {
			select {
			case mesasge := <-messages:
				if err := conn.WriteJSON(mesasge); err != nil {
					continue // TODO: log something or return the error???
				}

			case <-closeCh:
				conn.Close()
				return

			default:
				message := new(entities.Message)
				if err := conn.ReadJSON(message); err != nil {
					continue // TODO: log something or return the error???
				}
			}
		}
	}()

	return messages, closeCh, nil
}
