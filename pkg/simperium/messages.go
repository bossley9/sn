package simperium

import (
	l "git.sr.ht/~bossley9/sn/pkg/logger"

	"github.com/gorilla/websocket"
)

func writeMessage(conn *websocket.Conn, messageType int, message string) error {
	if err := conn.WriteMessage(messageType, []byte(message)); err != nil {
		return err
	}
	l.PrintDebug("\n" + "W " + message + "\n")

	return nil
}

func readMessage(conn *websocket.Conn) (int, string, error) {
	mtype, raw, err := conn.ReadMessage()
	if err != nil {
		return 0, "", err
	}
	message := string(raw)
	l.PrintDebug("\n" + "R " + message + "\n")

	return mtype, message, nil
}

func (client *Client) ReadMessage() (string, error) {
	_, message, err := readMessage(client.connection)
	if err != nil {
		return "", err
	}
	return message, nil
}
