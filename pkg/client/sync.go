package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"

	s "git.sr.ht/~bossley9/sn/pkg/simperium"
)

// sync client notes
func (client *client) Sync() error {
	if len(client.cache.CurrentVersion) == 0 {
		// initial sync
		fmt.Println("\tno change version found in cache. Making initial sync...")
		if err := client.initSync(); err != nil {
			return err
		}
	} else {
		// update sync
		fmt.Println("\tsyncing from version " + client.cache.CurrentVersion + "...")
		if err := client.updateSync(); err != nil {
			fmt.Println(err)
			fmt.Println("\tunable to update. Making initial sync...")
			if err := client.initSync(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (client *client) initSync() error {
	noteSummaries := []s.EntitySummary[Note]{}
	maxParallelNotes := 20

	isFirstBatch := true
	mark := ""
	version := ""

	for len(mark) > 0 || isFirstBatch {
		if len(mark) > 0 {
			fmt.Println("\t\tfetching batch " + mark + "...")
		} else {
			fmt.Println("\t\tfetching unmarked batch...")
		}

		if err := client.simp.WriteIndexMessage(0, true, mark, "", maxParallelNotes); err != nil {
			return err
		}
		message, err := client.simp.ReadMessage()
		if err != nil {
			return err
		}
		var indexRes s.IndexMessageResponse[Note]
		if err := json.Unmarshal([]byte(message[4:]), &indexRes); err != nil {
			return err
		}

		noteSummaries = append(noteSummaries, indexRes.Entities...)
		mark = indexRes.Mark

		if isFirstBatch {
			version = indexRes.CurrentVersion
			isFirstBatch = false
		}
	}

	for _, summary := range noteSummaries {
		if err := client.writeNoteSummary(&summary); err != nil {
			fmt.Println(err)
			continue
		}

		if err := client.saveNote(summary.ID, summary.Version); err != nil {
			fmt.Println(err)
			continue
		}
	}

	if err := client.setCurrentVersion(version); err != nil {
		return err
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
			fmt.Println("\t\tunable to retrieve file for entity " + change.EntityID + ". Skipping...")
			continue
		}
		content, err := os.ReadFile(filename)
		if err != nil {
			fmt.Println(err)
			fmt.Println("\t\tunable to read file " + filename + ". Skipping...")
			continue
		}

		fmt.Println("\tapplying change " + change.ChangeVersion + " to " + filename + "...")
		diff := change.Values.Content

		result := diff.Apply(string(content))

		if err := os.WriteFile(filename, []byte(result), 0600); err != nil {
			fmt.Println(err)
			fmt.Println("\t\tunable to update file " + filename + ". Skipping...")
			continue
		}

		fmt.Println("\tupdating change version from " + client.cache.CurrentVersion + " to " + change.ChangeVersion + "...")
		client.cache.CurrentVersion = change.ChangeVersion
		client.cache.Notes[change.EntityID] = change.EndVersion
		if err := client.writeCache(); err != nil {
			fmt.Println("\t\tunable to update cache. Skipping...")
			continue
		}
	}

	return nil
}
