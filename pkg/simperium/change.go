package simperium

import (
	"encoding/json"
	"strconv"
	"time"

	j "git.sr.ht/~bossley9/sn/pkg/jsondiff"

	"github.com/google/uuid"
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
	Error         int       `json:"error,omitempty"`
}

type UploadChange[T any] struct {
	SourceVersion int    `json:"sv"`
	EntityID      string `json:"id"`
	Operation     string `json:"o"`
	Values        T      `json:"v"`
	ChangeID      string `json:"ccid"`
}

type UploadDiff struct {
	Content          j.StringJSONDiff `json:"content"`
	ModificationDate j.Int64JSONDiff  `json:"modificationDate"`
}

func (client *Client) WriteChangeMessage(channel int, changeVersion string, entityVersion int, entityID string, operation string, textDiff string) (string, error) {
	ccid := uuid.New().String()
	contentDiff := UploadDiff{
		Content: j.StringJSONDiff{
			Operation: "d",
			Value:     textDiff,
		},
		ModificationDate: j.Int64JSONDiff{
			Operation: "r",
			Value:     time.Now().Unix(),
		},
	}
	change := UploadChange[UploadDiff]{
		SourceVersion: entityVersion,
		EntityID:      entityID,
		Operation:     operation,
		Values:        contentDiff,
		ChangeID:      ccid,
	}

	changeMsg, err := json.Marshal(change)
	if err != nil {
		return "", err
	}

	message := strconv.Itoa(channel) + ":c:" + string(changeMsg)
	if err := writeMessage(client.connection, websocket.TextMessage, message); err != nil {
		return "", err
	}
	return ccid, nil
}
