package client

import (
	"encoding/json"
	"fmt"
	"strings"

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
	maxNumNotes := 1
	if err := client.simp.WriteIndexMessage(0, false, "", "", maxNumNotes); err != nil {
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
	if err := WriteCache(client.cache); err != nil {
		return err
	}

	for _, entity := range indexRes.Entities {
		if err := client.simp.WriteEntityMessage(0, entity.ID, entity.Version); err != nil {
			fmt.Println(err)
			continue
		}
		entityMessage, err := client.simp.ReadMessage()
		if err != nil {
			return err
		}

		// remove first response line to parse data
		entityLines := strings.Split(entityMessage, "\n")
		_, entityLines = entityLines[0], entityLines[1:]
		entityMessage = strings.Join(entityLines, "\n")

		var noteRes s.EntityRes[Note]
		if err := json.Unmarshal([]byte(entityMessage), &noteRes); err != nil {
			return err
		}

		note := noteRes.Data

		note.ID = entity.ID
		note.Version = entity.Version

		client.writeNote(&note)
	}

	return nil
}
