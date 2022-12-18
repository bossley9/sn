package client

import (
	"encoding/json"
	"errors"
	"os"

	l "git.sr.ht/~bossley9/sn/pkg/logger"
	s "git.sr.ht/~bossley9/sn/pkg/simperium"
)

// sync client notes
func (client *Client) Sync() error {
	changeVersion := client.storage.ChangeVersion
	if len(changeVersion) == 0 {
		l.PrintWarning("Change version not found. Making fresh sync...")
		return client.RefetchSync()
	}

	l.PrintInfo("syncing from version " + changeVersion + "... ")
	if err := client.updateSync(); err != nil {
		l.PrintError(err)
		l.PrintInfo("\n")
		l.PrintWarning("Falling back to fresh sync... ")
		if err := client.RefetchSync(); err != nil {
			return err
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
			l.PrintInfo("\nFetching batch " + mark + "... ")
		} else {
			l.PrintInfo("\nFetching batch... ")
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
			l.PrintWarning("Warning: ")
			l.PrintWarning(err)
			l.PrintPlain("\n")
			continue
		}
	}

	client.storage.ChangeVersion = version
	return client.storage.writeChanges()
}

func (client *Client) updateSync() error {
	// force exit if local changes may conflict
	diffs := client.GetLocalDiffs()
	l.PrintPlain("\n")
	if len(diffs) > 0 {
		for noteID, diff := range diffs {
			noteCache, err := client.getCachedNote(noteID)
			if err != nil {
				l.PrintWarning("Unable to read local file with id " + noteID + ". Continuing...\n")
				continue
			}
			content, err := client.readVersionNote(noteID)
			if err != nil {
				l.PrintWarning("Unable to read local file with id " + noteID + ". Continuing...\n")
				continue
			}

			diff.PrettyPrint(noteCache.Name, content)
		}
		l.PrintWarning("Local diffs found. Please upload changes before syncing.\n")
		l.PrintPlain("\n")
		os.Exit(0)
	}

	changeVersion := client.storage.ChangeVersion
	if err := client.simp.WriteChangeVersionMessage(0, changeVersion); err != nil {
		return err
	}
	message, err := client.simp.ReadMessage()
	if err != nil {
		return err
	}

	if message == "0:cv:?" {
		return errors.New("change version does not exist.")
	} else if message == "0:c:[]" {
		l.PrintInfo("already up to date! ")
		return nil
	}

	changes, err := parseNoteChangeMessage(message)
	if err != nil {
		return err
	}

	// applying changes
	for _, change := range changes {
		client.applyChange(&change)
	}

	return client.storage.writeChanges()
}
