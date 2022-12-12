package client

import (
	"fmt"
	"os"
	"strconv"

	f "git.sr.ht/~bossley9/sn/pkg/fileio"
	j "git.sr.ht/~bossley9/sn/pkg/jsondiff"
)

// upload and sync local diffs with server
func (client *Client) Upload(diffs map[string]j.StringJSONDiff) error {
	for noteID, diff := range diffs {
		noteCache, err := client.getCachedNote(noteID)
		if err != nil {
			fmt.Println(err)
			fmt.Println("\tunable to find note with id " + noteID + " in cache. Continuing...")
			continue
		}

		ccid, err := client.simp.WriteChangeMessage(0, client.getCurrentVersion(), noteCache.Version, noteID, "M", diff.Value)
		if err != nil {
			fmt.Println(err)
			fmt.Println("\tunable to upload changes to " + noteCache.Name + ". Continuing...")
			continue
		}
		message, err := client.simp.ReadMessage()
		if err != nil {
			fmt.Println("\tunable to upload changes to " + noteCache.Name + ". Continuing...")
			continue
		}

		changes, err := parseNoteChangeMessage(message)
		if err != nil {
			return err
		}
		change := changes[0]
		if ccid != change.ChangeIDs[0] || change.Error > 0 {
			fmt.Println("\tunable to upload changes to " + noteCache.Name + " (error " + strconv.Itoa(change.Error) + "). Continuing...")
			continue
		}

		fmt.Println("\tchange successful.")

		fmt.Println("\tapplying changes...")
		// since changes are already applied locally, just copy over
		filename := client.getFileName(noteCache.Name)
		raw, err := os.ReadFile(filename)
		if err != nil {
			fmt.Println("\tunable to open note " + noteCache.Name + ". Continuing...")
			continue
		}
		vFilename := client.getVersionFileName(noteCache.Name)
		if err := os.WriteFile(vFilename, raw, f.RW); err != nil {
			fmt.Println("\tunable to writing version changes for note " + noteCache.Name + ". Continuing...")
			continue
		}

		if err := client.setNoteVersion(noteID, change.EndVersion); err != nil {
			fmt.Println("\tunable to update note version to " + strconv.Itoa(change.EndVersion) + ". Continuing...")
			continue
		}

		if err := client.setCurrentVersion(change.ChangeVersion); err != nil {
			fmt.Println("\tunable to update change version to " + change.ChangeVersion + ". Continuing...")
			continue
		}

		fmt.Println("\tchange applied to " + noteCache.Name + ".")
	}

	return client.writeCache()
}
