package simperium

import (
	"context"

	l "github.com/bossley9/sn/pkg/logger"

	"nhooyr.io/websocket"
)

func writeMessage(ctx context.Context, conn *websocket.Conn, messageType websocket.MessageType, message string) error {
	if err := conn.Write(ctx, messageType, []byte(message)); err != nil {
		return err
	}
	l.PrintDebug("\n" + "W " + message + "\n")

	return nil
}

func readMessage(ctx context.Context, conn *websocket.Conn) (websocket.MessageType, string, error) {
	mtype, raw, err := conn.Read(ctx)
	if err != nil {
		return 0, "", err
	}
	message := string(raw)
	l.PrintDebug("\n" + "R " + message + "\n")

	return mtype, message, nil
}

func (client *Client[DT]) ReadMessage(ctx context.Context) (string, error) {
	_, message, err := readMessage(ctx, client.connection)
	if err != nil {
		return "", err
	}
	return message, nil
}
