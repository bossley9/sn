package simperium

import (
	"encoding/json"
	"errors"
)

type AuthorizeResponse struct {
	AccessToken string `json:"access_token"`
	UserID      string `json:"userid"`
	Username    string `json:"username"`
}

// https://simperium.com/docs/reference/http/#auth
func (client *Client) Authorize(username string, password string) (string, error) {
	url := "https://auth.simperium.com/1/" + client.appID + "/authorize/"

	params := map[string]any{
		"username": username,
		"password": password,
	}

	headers := map[string]string{
		"X-Simperium-API-Key": client.apiKey,
	}

	authRes, err := fetch(url, "POST", params, headers)
	if err != nil {
		return "", err
	}

	var result AuthorizeResponse
	if err := json.Unmarshal(authRes, &result); err != nil {
		return "", errors.New(string(authRes))
	}

	return result.AccessToken, nil
}
