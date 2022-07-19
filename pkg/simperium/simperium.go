package simperium

import "github.com/gorilla/websocket"

type Client struct {
	appID         string
	apiKey        string
	apiVersion    string
	clientID      string
	clientName    string
	clientVersion string
	connection    *websocket.Conn
}

func NewClient(appID string, apiKey string, apiVer string, clientID string, clientName string, clientVer string) *Client {
	c := Client{
		appID:         appID,
		apiKey:        apiKey,
		apiVersion:    apiVer,
		clientID:      clientID,
		clientName:    clientName,
		clientVersion: clientVer,
	}
	return &c
}
