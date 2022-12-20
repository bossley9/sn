package simperium

import (
	"context"
	"strconv"

	"nhooyr.io/websocket"
)

func (client *Client[DT]) WriteEntityMessage(ctx context.Context, channel int, entityID string, entityVersion int) error {
	message := strconv.Itoa(channel) + ":e:" + entityID + "." + strconv.Itoa(entityVersion)
	if err := writeMessage(ctx, client.connection, websocket.MessageText, message); err != nil {
		return err
	}
	return nil
}
