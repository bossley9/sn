package simperium

import (
	"context"
	"strconv"

	"nhooyr.io/websocket"
)

type IndexMessageResponse[T any] struct {
	CurrentVersion string             `json:"current"`
	Entities       []EntitySummary[T] `json:"index"`
	Mark           string             `json:"mark"`
}

type EntitySummary[T any] struct {
	ID      string `json:"id"`
	Version int    `json:"v"`
	Data    T      `json:"d,omitempty"`
}

func (client *Client) WriteIndexMessage(ctx context.Context, channel int, returnData bool, offset string, mark string, limit int) error {
	message := strconv.Itoa(channel) + ":i:"

	if returnData {
		message = message + "1"
	}
	message = message + ":"

	if len(offset) > 0 {
		message = message + offset
	}
	message = message + ":"

	if len(mark) > 0 {
		message = message + mark
	}
	message = message + ":"

	// limit is 0-indexed
	message = message + strconv.Itoa(limit-1)

	if err := writeMessage(ctx, client.connection, websocket.MessageText, message); err != nil {
		return err
	}

	return nil
}
