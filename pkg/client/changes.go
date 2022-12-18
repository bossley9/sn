package client

import (
	"encoding/json"
	"os"

	j "git.sr.ht/~bossley9/sn/pkg/jsondiff"
	l "git.sr.ht/~bossley9/sn/pkg/logger"
	s "git.sr.ht/~bossley9/sn/pkg/simperium"
)

func parseNoteChangeMessage(message string) (s.ChangeVersionResponse[NoteDiff], error) {
	response := message[4:]
	var changes s.ChangeVersionResponse[NoteDiff]
	if err := json.Unmarshal([]byte(response), &changes); err != nil {
		return nil, err
	}
	return changes, nil
}

// given any change, applies that change to the specified note
func (client *Client) applyChange(change *s.Change[NoteDiff]) {
	noteID := change.EntityID

	if change.Values.CreationDate.Operation == j.OP_INSERT {
		// note creation
		l.PrintInfo("Creating note " + noteID + "... ")
		if err := client.applyCreationChange(change); err != nil {
			l.PrintError(err)
			l.PrintWarning("\nUnable to create note " + noteID + ". Skipping...")
		}
	} else if change.Operation == j.OP_DELETE {
		// note deletion
		l.PrintInfo("Deleting note " + noteID + "... ")
		if err := client.applyDeletionChange(change); err != nil {
			l.PrintError(err)
			l.PrintWarning("\nUnable to delete note " + noteID + ". Skipping...")
		}
	} else if len(change.Values.Content.Value) > 0 {
		// note update
		l.PrintInfo("Updating note " + noteID + "... ")
		if err := client.applyUpdateChange(change); err != nil {
			l.PrintError(err)
			l.PrintWarning("\nUnable to update note " + noteID + ". Skipping...")
		}
	} else {
		// silently ignore other changes such as pinning or toggling markdown
	}

	l.PrintInfo("\nUpdating change version to " + change.ChangeVersion + "... \n")
	client.storage.ChangeVersion = change.ChangeVersion
}

// given an update change, applies that change to the specified note
func (client *Client) applyUpdateChange(change *s.Change[NoteDiff]) error {
	noteID := change.EntityID

	content, err := client.readNote(noteID)
	if err != nil {
		return err
	}

	l.PrintInfo("Applying change " + change.ChangeVersion + " to note " + noteID + "... ")
	result := change.Values.Content.Apply(string(content))

	l.PrintInfo("writing changes... ")
	noteSummary := NoteSummary{
		ID:      noteID,
		Version: change.EndVersion,
		Content: result,
	}
	if err := client.writeNote(&noteSummary); err != nil {
		return err
	}

	return nil
}

// given a creation change, applies that change to the specified note
func (client *Client) applyCreationChange(change *s.Change[NoteDiff]) error {
	summary := NoteSummary{
		ID:      change.EntityID,
		Version: change.EndVersion,
		Content: change.Values.Content.Value,
	}
	return client.writeNote(&summary)
}

// given a deletion change, deletes the specified note
func (client *Client) applyDeletionChange(change *s.Change[NoteDiff]) error {
	noteID := change.EntityID
	noteCache, err := client.getCachedNote(noteID)
	if err != nil {
		return err
	}

	// remove file
	filename := client.getFileName(noteCache.Name)
	if err := os.Remove(filename); err != nil {
		return err
	}
	vFilename := client.getVersionFileName(noteCache.Name)
	if err := os.Remove(vFilename); err != nil {
		return err
	}

	// remove from cache
	delete(client.cache.Notes, noteID)

	return nil
}
