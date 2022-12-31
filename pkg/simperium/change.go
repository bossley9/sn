package simperium

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/google/uuid"
	"nhooyr.io/websocket"

	j "git.sr.ht/~bossley9/sn/pkg/jsondiff"
)

type Change[T interface{}] struct {
	ClientID      string    `json:"clientid,omitempty"`
	ChangeVersion string    `json:"cv,omitempty"`
	EndVersion    int       `json:"ev,omitempty"`
	SourceVersion int       `json:"sv,omitempty"`
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
	var changeMessage string

	if operation == j.OP_DELETE {
		change := Change[string]{
			EntityID:  entityID,
			Operation: operation,
			Values:    "",
			ChangeID:  ccid,
		}
		changeMsg, err := json.Marshal(change)
		if err != nil {
			return "", err
		}
		changeMessage = string(changeMsg)

	} else {
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
		changeMessage = string(changeMsg)

	}

	message := strconv.Itoa(channel) + ":c:" + changeMessage
	if err := writeMessage(ctx, client.connection, websocket.MessageText, message); err != nil {
		return "", err
	}
	return ccid, nil
}
