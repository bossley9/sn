package simperium

import (
	"strconv"

	"github.com/gorilla/websocket"
)

type EntityRes[T any] struct {
	Data T `json:"data"`
}

func (client *Client) WriteEntityMessage(channel int, entityID string, entityVersion int) error {
	message := strconv.Itoa(channel) + ":e:" + entityID + "." + strconv.Itoa(entityVersion)
	if err := writeMessage(client.connection, websocket.TextMessage, message); err != nil {
		return err
	}
	return nil
}
