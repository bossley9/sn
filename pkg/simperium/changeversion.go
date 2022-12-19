package simperium

import (
	"context"
	"strconv"

	"nhooyr.io/websocket"
)

func (client *Client) WriteChangeVersionMessage(ctx context.Context, channel int, changeVersion string) error {
	message := strconv.Itoa(channel) + ":cv:" + changeVersion
	if err := writeMessage(ctx, client.connection, websocket.MessageText, message); err != nil {
		return err
	}
	return nil
}
