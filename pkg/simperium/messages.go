package simperium

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gorilla/websocket"
)

// TODO remove debug
func writeMessage(conn *websocket.Conn, messageType int, message string, debug bool) error {
	if err := conn.WriteMessage(messageType, []byte(message)); err != nil {
		return err
	}
	if debug {
		fmt.Println("w " + message)
	}

	return nil
}

// TODO remove debug
func readMessage(conn *websocket.Conn, debug bool) (int, string, error) {
	mtype, raw, err := conn.ReadMessage()
	if err != nil {
		return 0, "", err
	}
	message := string(raw)
	if debug {
		fmt.Println("r " + message)
	}

	return mtype, message, nil
}

func (client *Client) ReadMsg() (string, error) {
	_, message, err := readMessage(client.connection, true)
	if err != nil {
		return "", err
	}
	return message, nil
}

type InitMessage struct {
	ClientID   string `json:"clientid"`
	API        string `json:"api"`
	Token      string `json:"token"`
	AppID      string `json:"app_id"`
	BucketName string `json:"name"`
	Library    string `json:"library"`
	Version    string `json:"version"`
}

func (client *Client) WriteInitMsg(channel int, token string, bucketName string) error {
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

	if err := writeMessage(client.connection, websocket.TextMessage, message, true); err != nil {
		return err
	}

	return nil
}

func (client *Client) WriteIndexMsg(channel int, returnData bool, offset string, mark string, limit int) error {
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

	message = message + strconv.Itoa(limit)

	if err := writeMessage(client.connection, websocket.TextMessage, message, true); err != nil {
		return err
	}

	return nil
}
