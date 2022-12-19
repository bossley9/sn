package simperium

import (
	"strconv"

	"github.com/gorilla/websocket"
)

func (client *Client) WriteChangeVersionMessage(channel int, changeVersion string) error {
	message := strconv.Itoa(channel) + ":cv:" + changeVersion
	if err := writeMessage(client.connection, websocket.TextMessage, message); err != nil {
		return err
	}
	return nil
}
