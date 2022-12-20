package simperium

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/google/uuid"
	"nhooyr.io/websocket"
)

type Change[T interface{}] struct {
	ClientID      string    `json:"clientid"`
	ChangeVersion string    `json:"cv"`
	EndVersion    int       `json:"ev"`
	SourceVersion int       `json:"sv"`
	EntityID      string    `json:"id"`
	Operation     string    `json:"o"`
	Values        T         `json:"v"`
	ChangeIDs     []string  `json:"ccids"`
	Data          *struct{} `json:"d,omitempty"`
	Error         int       `json:"error,omitempty"`
}

type UploadChange[T interface{}] struct {
	SourceVersion int       `json:"sv"`
	EntityID      string    `json:"id"`
	Operation     string    `json:"o"`
	Values        T         `json:"v"`
	ChangeID      string    `json:"ccid"`
	Data          *struct{} `json:"d,omitempty"`
}

func (client *Client[DT]) WriteChangeMessage(ctx context.Context, channel int, changeVersion string, entityVersion int, entityID string, operation string, diff DT) (string, error) {
	ccid := uuid.New().String()
	change := UploadChange[DT]{
		SourceVersion: entityVersion,
		EntityID:      entityID,
		Operation:     operation,
		Values:        diff,
		ChangeID:      ccid,
	}

	changeMsg, err := json.Marshal(change)
	if err != nil {
		return "", err
	}

	message := strconv.Itoa(channel) + ":c:" + string(changeMsg)
	if err := writeMessage(ctx, client.connection, websocket.MessageText, message); err != nil {
		return "", err
	}
	return ccid, nil
}
