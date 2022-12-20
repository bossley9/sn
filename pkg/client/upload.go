package client

import (
	"context"
	"os"
	"strconv"

	f "git.sr.ht/~bossley9/sn/pkg/fileio"
	j "git.sr.ht/~bossley9/sn/pkg/jsondiff"
	l "git.sr.ht/~bossley9/sn/pkg/logger"
)

// upload and sync local diffs with server
func (client *Client) Upload(ctx context.Context, diffs map[NoteID]UploadNoteDiff) error {
	for noteID, diff := range diffs {
		note, ok := client.storage.Notes[noteID]
		if !ok {
			l.PrintWarning("\nUnable to find note with id " + noteID + " in cache. Continuing...\n")
			continue
		}

		changeVersion := client.storage.ChangeVersion
		ccid, err := client.simp.WriteChangeMessage(ctx, 0, changeVersion, note.Version, string(noteID), j.OP_MODIFY, diff.Values)
		if err != nil {
			l.PrintError(err)
			l.PrintWarning("\nUnable to upload changes to " + note.Name + ". Continuing...\n")
			continue
		}
		message, err := client.simp.ReadMessage(ctx)
		if err != nil {
			l.PrintError(err)
			l.PrintWarning("\nUnable to upload changes to " + note.Name + ". Continuing...\n")
			continue
		}

		changes, err := parseNoteChangeMessage(message)
		if err != nil {
			return err
		}
		change := changes[0]
		if ccid != change.ChangeIDs[0] || change.Error > 0 {
			// https://simperium.com/docs/websocket/#change-c
			var errorMessage string
			switch change.Error {
			case 440:
				errorMessage = "invalid diff. Try removing non-ascii or non-traditional characters"
			default:
				errorMessage = "internal server error"
			}
			l.PrintError("error " + strconv.Itoa(change.Error) + ": " + errorMessage + ".")
			l.PrintWarning("\nUnable to upload changes to " + note.Name + ". Continuing...\n")
			continue
		}

		// applying changes
		// since changes are already applied locally, just copy over
		filename := client.getFileName(note.Name)
		raw, err := os.ReadFile(filename)
		if err != nil {
			l.PrintWarning("Unable to open note " + note.Name + ". Continuing...\n")
			continue
		}
		vFilename := client.getVersionFileName(note.Name)
		if err := os.WriteFile(vFilename, raw, f.RW); err != nil {
			l.PrintWarning("Unable to writing version changes for note " + note.Name + ". Continuing...\n")
			continue
		}

		note.Version = change.EndVersion

		client.storage.Notes[noteID] = note
		client.storage.ChangeVersion = change.ChangeVersion

		l.PrintInfo("Change applied to " + note.Name + ".\n")
	}

	return client.storage.writeChanges()
}
