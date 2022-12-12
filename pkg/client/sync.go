package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	s "git.sr.ht/~bossley9/sn/pkg/simperium"
)

// sync client notes
func (client *Client) Sync() error {
	currentVersion := client.getCurrentVersion()
	if len(currentVersion) == 0 {
		fmt.Println("\tno change version found in cache. Making initial sync...")
		if err := client.RefetchSync(); err != nil {
			return err
		}
	} else {
		fmt.Println("\tsyncing from version " + currentVersion + "...")
		if err := client.updateSync(); err != nil {
			fmt.Println(err)
			fmt.Println("\tunable to update. Falling back to initial sync...")
			if err := client.RefetchSync(); err != nil {
				return err
			}
		}
	}

	return nil
}

// initial sync to load (or reload) all notes
func (client *Client) RefetchSync() error {
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

func (client *Client) updateSync() error {
	// force exit if local changes may conflict
	diffs, err := client.GetLocalDiffs()
	if err != nil {
		return err
	}
	if len(diffs) > 0 {
		for _, diff := range diffs {
			fmt.Println("\t\t" + diff.Value)
		}
		log.Fatal("local diffs found. Please upload before syncing. Exiting.")
	}

	if err := client.simp.WriteChangeVersionMessage(0, client.getCurrentVersion()); err != nil {
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

	changes, err := parseNoteChangeMessage(message)
	if err != nil {
		return err
	}

	fmt.Println("\tapplying changes...")
	for _, change := range changes {
		client.applyChange(&change)
	}

	return nil
}
