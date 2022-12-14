package client

import (
	"os"
	"strconv"

	f "git.sr.ht/~bossley9/sn/pkg/fileio"
	j "git.sr.ht/~bossley9/sn/pkg/jsondiff"
	l "git.sr.ht/~bossley9/sn/pkg/logger"
)

// upload and sync local diffs with server
func (client *Client) Upload(diffs map[string]j.StringJSONDiff) error {
	for noteID, diff := range diffs {
		noteCache, err := client.getCachedNote(noteID)
		if err != nil {
			l.PrintError(err)
			l.PrintWarning("\nUnable to find note with id " + noteID + " in cache. Continuing...\n")
			continue
		}

		ccid, err := client.simp.WriteChangeMessage(0, client.getCurrentVersion(), noteCache.Version, noteID, "M", diff.Value)
		if err != nil {
			l.PrintError(err)
			l.PrintWarning("\nUnable to upload changes to " + noteCache.Name + ". Continuing...\n")
			continue
		}
		message, err := client.simp.ReadMessage()
		if err != nil {
			l.PrintError(err)
			l.PrintWarning("\nUnable to upload changes to " + noteCache.Name + ". Continuing...\n")
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
			l.PrintWarning("\nUnable to upload changes to " + noteCache.Name + ". Continuing...\n")
			continue
		}

		// applying changes
		// since changes are already applied locally, just copy over
		filename := client.getFileName(noteCache.Name)
		raw, err := os.ReadFile(filename)
		if err != nil {
			l.PrintWarning("Unable to open note " + noteCache.Name + ". Continuing...\n")
			continue
		}
		vFilename := client.getVersionFileName(noteCache.Name)
		if err := os.WriteFile(vFilename, raw, f.RW); err != nil {
			l.PrintWarning("Unable to writing version changes for note " + noteCache.Name + ". Continuing...\n")
			continue
		}

		if err := client.setNoteVersion(noteID, change.EndVersion); err != nil {
			l.PrintWarning("Unable to update note version to " + strconv.Itoa(change.EndVersion) + ". Continuing...\n")
			continue
		}

		if err := client.setCurrentVersion(change.ChangeVersion); err != nil {
			l.PrintWarning("Unable to update change version to " + change.ChangeVersion + ". Continuing...\n")
			continue
		}

		l.PrintInfo("Change applied to " + noteCache.Name + ".\n")
	}

	return client.writeCache()
}
