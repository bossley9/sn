package simperium

import (
	"context"
	"strconv"

	"nhooyr.io/websocket"
)

type EntityRes[T any] struct {
	Data T `json:"data"`
}

func (client *Client) WriteEntityMessage(ctx context.Context, channel int, entityID string, entityVersion int) error {
	message := strconv.Itoa(channel) + ":e:" + entityID + "." + strconv.Itoa(entityVersion)
	if err := writeMessage(ctx, client.connection, websocket.MessageText, message); err != nil {
		return err
	}
	return nil
}
