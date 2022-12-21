package simperium

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/google/uuid"
	"nhooyr.io/websocket"
)

type Change[T interface{}] struct {
	ClientID      string    `json:"clientid,omitempty"`
	ChangeVersion string    `json:"cv,omitempty"`
	EndVersion    int       `json:"ev,omitempty"`
	SourceVersion int       `json:"sv"`
	EntityID      string    `json:"id"`
	Operation     string    `json:"o"`
	Values        T         `json:"v,omitempty"`
	ChangeID      string    `json:"ccid,omitempty"`  // download changes only
	ChangeIDs     []string  `json:"ccids,omitempty"` // upload changes only
	Data          *struct{} `json:"d,omitempty"`
	Error         int       `json:"error,omitempty"`
}

func (client *Client[DT]) WriteChangeMessage(ctx context.Context, channel int, changeVersion string, entityVersion int, entityID string, operation string, diff DT) (string, error) {
	ccid := uuid.New().String()
	change := Change[DT]{
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
