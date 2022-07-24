package simperium

import (
	"encoding/json"
	"strconv"

	"github.com/gorilla/websocket"
)

type Change[T any] struct {
	ClientID      string    `json:"clientid"`
	ChangeVersion string    `json:"cv"`
	EndVersion    int       `json:"ev"`
	SourceVersion int       `json:"sv"`
	EntityID      string    `json:"id"`
	Operation     string    `json:"o"`
	Values        T         `json:"v"`
	ChangeIDs     []string  `json:"ccids"`
	Data          *struct{} `json:"d,omitempty"`
}

func (client *Client) WriteChangeMessage(channel int, changeVersion string, currentVersion int, entityID string, operation string, values *struct{}) error {
	change := Change[*struct{}]{
		ClientID:      client.clientID,
		ChangeVersion: changeVersion,
		EndVersion:    currentVersion + 1,
		SourceVersion: currentVersion,
		EntityID:      entityID,
		Operation:     operation,
		Values:        values,
		ChangeIDs:     []string{}, // TODO generate ID
		Data:          nil,
	}

	changeMsg, err := json.Marshal(change)
	if err != nil {
		return err
	}

	message := strconv.Itoa(channel) + ":c:" + string(changeMsg)
	if err := writeMessage(client.connection, websocket.TextMessage, message); err != nil {
		return err
	}
	return nil
}
