package client

import (
	"encoding/json"
	"errors"
	"fmt"

	s "git.sr.ht/~bossley9/sn/pkg/simperium"
)

// sync client notes
func (client *client) Sync() error {
	if len(client.cache.CurrentVersion) == 0 {
		fmt.Println("\tno change version found in cache. Making initial sync...")
		if err := client.initSync(); err != nil {
			return err
		}
	} else {
		fmt.Println("\tsyncing from version " + client.cache.CurrentVersion + "...")
		if err := client.updateSync(); err != nil {
			fmt.Println(err)
			fmt.Println("\tunable to update. Falling back to initial sync...")
			if err := client.initSync(); err != nil {
				return err
			}
		}
	}

	return nil
}

// initial sync to load (or reload) all notes
func (client *client) initSync() error {
	noteSummaries := make([]s.EntitySummary[Note], 0)
	maxParallelNotes := 20

	isFirstBatch := true
	mark := ""
	version := ""

	// batch fetch notes
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

	// write notes
	for _, summary := range noteSummaries {
		// reformat data
		noteSummary := NoteSummary{
			ID:      summary.ID,
			Version: summary.Version,
			Content: summary.Data.Content,
		}

		if err := client.writeNote(&noteSummary); err != nil {
			fmt.Println(err)
			continue
		}
	}

	// update current version
	if err := client.setCurrentVersion(version); err != nil {
		return err
	}

	return nil
}

func (client *client) updateSync() error {
	if err := client.simp.WriteChangeVersionMessage(0, client.cache.CurrentVersion); err != nil {
		return err
	}
	message, err := client.simp.ReadMessage()
	if err != nil {
		return err
	}

	if message == "0:cv:?" {
		return errors.New("change version does not exist for bucket")

	} else if message == "0:c:[]" {
		fmt.Println("\tclient is already up to date!")
		return nil
	}

	response := message[4:]
	var changes s.ChangeVersionResponse[NoteDiff]
	if err := json.Unmarshal([]byte(response), &changes); err != nil {
		return err
	}

	fmt.Println("\tapplying changes...")
	for _, change := range changes {
		client.applyChange(&change)
	}

	return nil
}

func (client *client) applyChange(change *s.Change[NoteDiff]) {
	noteID := change.EntityID

	// read note
	content, err := client.readNote(noteID)
	if err != nil {
		fmt.Println(err)
		fmt.Println("\t\tunable to retrieve note " + noteID + ". Skipping...")
		return
	}

	// apply diff
	fmt.Println("\t\tapplying change " + change.ChangeVersion + " to note " + noteID + "...")
	result := change.Values.Content.Apply(string(content))

	// write note
	fmt.Println("\t\twriting changes...")
	noteSummary := NoteSummary{
		ID:      noteID,
		Version: change.EndVersion,
		Content: result,
	}
	if err := client.writeNote(&noteSummary); err != nil {
		fmt.Println("\t\tunable to update note " + noteSummary.ID + ". Skipping...")
		return
	}

	// update change version
	fmt.Println("\t\tupdating change version from " + client.cache.CurrentVersion + " to " + change.ChangeVersion + "...")
	if err := client.setCurrentVersion(change.ChangeVersion); err != nil {
		fmt.Println("\t\tunable to set current version. Skipping...")
		return
	}
}
