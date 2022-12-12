package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	s "git.sr.ht/~bossley9/sn/pkg/simperium"
)

const red = "\033[0;31m"
const yellow = "\033[0;33m"
const cyan = "\033[0;36m"
const none = "\033[0m"

// sync client notes
func (client *Client) Sync() error {
	currentVersion := client.getCurrentVersion()
	if len(currentVersion) == 0 {
		fmt.Print(yellow)
		fmt.Print("no version found in cache. Making fresh sync...")
		fmt.Print(cyan)
		if err := client.RefetchSync(); err != nil {
			return err
		}
	} else {
		fmt.Print("syncing from version " + currentVersion + "... ")
		if err := client.updateSync(); err != nil {
			fmt.Print(red)
			fmt.Println(err)
			fmt.Print(yellow)
			fmt.Print("Falling back to fresh sync... ")
			fmt.Print(cyan)
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
			fmt.Print(cyan)
			fmt.Print("\nFetching batch " + mark + "... ")
		} else {
			fmt.Print(cyan)
			fmt.Print("\nFetching batch... ")
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
	fmt.Println(none)
	if len(diffs) > 0 {
		for _, diff := range diffs {
			fmt.Println(diff.Value)
			// TODO display a richer diff for usability
		}
		fmt.Println(yellow)
		fmt.Println("Local diffs found. Please upload changes before syncing.")
		fmt.Println(none)
		os.Exit(0)
	}

	if err := client.simp.WriteChangeVersionMessage(0, client.getCurrentVersion()); err != nil {
		return err
	}
	message, err := client.simp.ReadMessage()
	if err != nil {
		return err
	}

	if message == "0:cv:?" {
		return errors.New("change version does not exist.")
	} else if message == "0:c:[]" {
		fmt.Print(cyan)
		fmt.Print("already up to date! ")
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
