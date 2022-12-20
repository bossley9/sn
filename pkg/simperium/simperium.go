package simperium

import "nhooyr.io/websocket"

type Client[DiffType interface{}] struct {
	appID         string
	apiKey        string
	apiVersion    string
	clientID      string
	clientName    string
	clientVersion string
	connection    *websocket.Conn
}

func NewClient[DT interface{}](appID string, apiKey string, apiVer string, clientID string, clientName string, clientVer string) *Client[DT] {
	return &Client[DT]{
		appID:         appID,
		apiKey:        apiKey,
		apiVersion:    apiVer,
		clientID:      clientID,
		clientName:    clientName,
		clientVersion: clientVer,
	}
}
