package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	s "git.sr.ht/~bossley9/sn/pkg/simperium"
)

// sync client with bucket
func (client *client) Sync() error {
	if len(client.cache.CurrentVersion) == 0 {
		// initial sync
		fmt.Println("\tno change version found. Making initial sync...")
		if err := client.initialSync(); err != nil {
			return err
		}
	} else {
		// update sync
		fmt.Println("\tsyncing from version " + client.cache.CurrentVersion + "...")
		if err := client.updateSync(); err != nil {
			fmt.Println(err)
			fmt.Println("\tunable to update. Making initial sync...")
			if err := client.initialSync(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (client *client) initialSync() error {
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
			fmt.Println(err)
			continue
		}

		// remove first response line to parse data
		entityLines := strings.Split(entityMessage, "\n")
		_, entityLines = entityLines[0], entityLines[1:]
		entityMessage = strings.Join(entityLines, "\n")

		var noteRes s.EntityRes[Note]
		if err := json.Unmarshal([]byte(entityMessage), &noteRes); err != nil {
			fmt.Println(err)
			continue
		}

		note := noteRes.Data

		note.ID = entity.ID
		note.Version = entity.Version

		if err := client.writeNote(&note); err != nil {
			fmt.Println(err)
			continue
		}
	}

	return nil
}

func (client *client) updateSync() error {
	channel := 0
	if err := client.simp.WriteChangeVersionMessage(channel, client.cache.CurrentVersion); err != nil {
		return err
	}
	message, err := client.simp.ReadMessage()
	if err != nil {
		return err
	}

	channelText := strconv.Itoa(channel)
	if message == channelText+":cv:?" {
		// change version does not exist for bucket
		return errors.New("change version does not exist for bucket")
	} else if message == channelText+":c:[]" {
		// client is up to date
		fmt.Println("\tclient is up to date!")
		return nil
	}

	response := message[4:]
	var changes s.ChangeVersionResponse[NoteDiff]
	if err := json.Unmarshal([]byte(response), &changes); err != nil {
		return err
	}

	fmt.Println("\tapplying changes...")
	for _, change := range changes {
		filename, err := client.getFileNameFromID(change.EntityID)
		if err != nil {
			fmt.Println(err)
			fmt.Println("\t\tunable to apply change to entity " + change.EntityID + ". Skipping...")
			continue
		}
		fmt.Println("\tapplying change " + change.ChangeVersion + " to " + filename + "...")
		// TODO
	}

	return nil
}
