package client

import (
	"encoding/json"
	"fmt"

	s "git.sr.ht/~bossley9/sn/pkg/simperium"
)

// sync client with bucket
func (client *client) Sync() error {
	// first sync
	fmt.Println("\tmaking first sync...")
	if err := client.doFirstSync(); err != nil {
		return err
	}

	return nil
}

func (client *client) doFirstSync() error {
	if err := client.simp.WriteIndexMessage(0, false, "", "", 10); err != nil {
		return err
	}
	message, err := client.simp.ReadMessage()
	if err != nil {
		return err
	}
	var indexRes s.IndexMessageResponse
	if err := json.Unmarshal([]byte(message[4:]), &indexRes); err != nil {
		return err
	}
	client.cache.CurrentVersion = indexRes.CurrentVersion
	WriteCache(client.cache)

	return nil
}
