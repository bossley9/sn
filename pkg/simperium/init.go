package simperium

import (
	"context"
	"encoding/json"
	"strconv"

	"nhooyr.io/websocket"
)

type InitMessage struct {
	ClientID   string `json:"clientid"`
	API        string `json:"api"`
	Token      string `json:"token"`
	AppID      string `json:"app_id"`
	BucketName string `json:"name"`
	Library    string `json:"library"`
	Version    string `json:"version"`
}

func (client *Client) WriteInitMessage(ctx context.Context, channel int, token string, bucketName string) error {
	messageJson := InitMessage{
		ClientID:   client.clientID,
		API:        client.apiVersion,
		Token:      token,
		AppID:      client.appID,
		BucketName: bucketName,
		Library:    client.clientName,
		Version:    client.clientVersion,
	}

	messageBytes, err := json.Marshal(messageJson)
	if err != nil {
		return err
	}

	message := strconv.Itoa(channel) + ":init:" + string(messageBytes)

	if err := writeMessage(ctx, client.connection, websocket.MessageText, message); err != nil {
		return err
	}

	return nil
}
