package simperium

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func fetch(url string, method string, params map[string]any, headers map[string]string) ([]byte, error) {
	bodyParams, err := json.Marshal(params)
	if err != nil {
		return []byte{}, err
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(bodyParams))
	if err != nil {
		return []byte{}, err
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Accept", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}

	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
