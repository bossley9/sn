package client

import (
	"context"
	"encoding/json"
	"errors"
	"os"

	l "git.sr.ht/~bossley9/sn/pkg/logger"
	s "git.sr.ht/~bossley9/sn/pkg/simperium"
)

// sync client notes
func (client *Client) Sync(ctx context.Context) error {
	changeVersion := client.storage.ChangeVersion
	if len(changeVersion) == 0 {
		l.PrintWarning("Change version not found. Making fresh sync...")
		return client.RefetchSync(ctx)
	}

	l.PrintInfo("syncing from version " + changeVersion + "... ")
	err := client.updateSync(ctx)
	if err != nil {
		l.PrintError(err)
		l.PrintWarning("\nFalling back to fresh sync... ")
		return client.RefetchSync(ctx)
	}

	return nil
}

// initial sync to load (or reload) all notes
func (client *Client) RefetchSync(ctx context.Context) error {
	noteEntities := []s.EntitySummary[NoteResponse]{}
	maxParallelNotes := 30

	mark := ""
	version := ""

	// first batch does not include mark
	isFirstBatch := true
	for len(mark) > 0 || isFirstBatch {
		if len(mark) > 0 {
			l.PrintInfo("\nFetching batch " + mark + "... ")
		} else {
			l.PrintInfo("\nFetching batch... ")
		}

		if err := client.simp.WriteIndexMessage(ctx, 0, true, mark, "", maxParallelNotes); err != nil {
			return err
		}
		message, err := client.simp.ReadMessage(ctx)
		if err != nil {
			return err
		}
		var indexRes s.IndexMessageResponse[NoteResponse]
		if err := json.Unmarshal([]byte(message[4:]), &indexRes); err != nil {
			return err
		}

		noteEntities = append(noteEntities, indexRes.Entities...)
		mark = indexRes.Mark

		if isFirstBatch {
			version = indexRes.CurrentVersion
			isFirstBatch = false
		}
	}

	for _, noteResponse := range noteEntities {
		noteID := NoteID(noteResponse.ID)
		content := noteResponse.Data.Content

		note := Note{
			Version: noteResponse.Version,
			Name:    GetNoteName(noteID, content),
		}

		if err := client.writeNote(noteID, &note, noteResponse.Data.Content); err != nil {
			l.PrintWarning("Warning: ")
			l.PrintWarning(err)
			l.PrintPlain("\n")
			continue
		}
	}

	client.storage.ChangeVersion = version
	return client.storage.writeChanges()
}

func (client *Client) updateSync(ctx context.Context) error {
	// force exit if local changes may conflict
	diffs := client.GetLocalDiffs()
	l.PrintPlain("\n")
	if len(diffs) > 0 {
		for noteID, diff := range diffs {
			note, ok := client.storage.Notes[noteID]
			if !ok {
				l.PrintWarning("Unable to read local file with id " + noteID + ". Continuing...\n")
				continue
			}
			content, err := client.readVersionNote(noteID)
			if err != nil {
				l.PrintWarning("Unable to read local file with id " + noteID + ". Continuing...\n")
				continue
			}

			diff.PrettyPrint(note.Name, content)
		}
		l.PrintWarning("Local diffs found. Please upload changes before syncing.\n")
		l.PrintPlain("\n")
		os.Exit(0)
	}

	changeVersion := client.storage.ChangeVersion
	if err := client.simp.WriteChangeVersionMessage(ctx, 0, changeVersion); err != nil {
		return err
	}
	message, err := client.simp.ReadMessage(ctx)
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
