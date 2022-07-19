package simperium

import (
	"time"

	"github.com/gorilla/websocket"
)

func (client *Client) ConnectToSocket() error {
	dialer := websocket.DefaultDialer
	dialer.EnableCompression = true // experimental
	dialer.HandshakeTimeout = 30 * time.Second

	socketUrl := "wss://api.simperium.com/sock/1/" + client.appID + "/websocket"
	conn, _, err := dialer.Dial(socketUrl, nil)
	if err != nil {
		return err
	}

	client.connection = conn
	return nil
}

func (client *Client) DisconnectSocket() error {
	err := client.connection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		return err
	}
	return nil
}
