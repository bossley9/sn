package simperium

import (
	"context"

	"nhooyr.io/websocket"
)

func (client *Client) ConnectToSocket(ctx context.Context) error {
	dialOpts := websocket.DialOptions{
		CompressionMode: websocket.CompressionContextTakeover,
	}

	socketUrl := "wss://api.simperium.com/sock/1/" + client.appID + "/websocket"
	conn, _, err := websocket.Dial(ctx, socketUrl, &dialOpts)
	if err != nil {
		return err
	}
	// raise initial read limit of 32768 (2^15) bytes to prevent rate limiting
	conn.SetReadLimit(524288) // 2^19 bytes

	client.connection = conn
	return nil
}

func (client *Client) DisconnectSocket() error {
	err := client.connection.Close(websocket.StatusNormalClosure, "")
	if err != nil {
		return err
	}
	return nil
}
