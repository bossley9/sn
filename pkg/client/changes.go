package client

import (
	"encoding/json"
	"errors"
	"os"

	j "git.sr.ht/~bossley9/sn/pkg/jsondiff"
	l "git.sr.ht/~bossley9/sn/pkg/logger"
)

func parseNoteChangeMessage(message string) ([]NoteChange, error) {
	response := message[4:]
	var changes []NoteChange
	if err := json.Unmarshal([]byte(response), &changes); err != nil {
		return nil, err
	}
	return changes, nil
}

// given any change, applies that change to the specified note
func (client *Client) applyChange(change *NoteChange) {
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
func (client *Client) applyUpdateChange(change *NoteChange) error {
	noteID := NoteID(change.EntityID)

	content, err := client.readNote(noteID)
	if err != nil {
		return err
	}

	l.PrintInfo("Applying change " + change.ChangeVersion + " to note " + string(noteID) + "... ")
	modifiedContent := change.Values.Content.Apply(content)

	l.PrintInfo("writing changes... ")
	note := Note{
		Version: change.EndVersion,
		Name:    client.GetNoteName(noteID, modifiedContent),
	}
	if err := client.writeNote(noteID, &note, modifiedContent); err != nil {
		return err
	}

	return nil
}

// given a creation change, applies that change to the specified note
func (client *Client) applyCreationChange(change *NoteChange) error {
	noteID := NoteID(change.EntityID)
	content := change.Values.Content.Value
	note := Note{
		Version: change.EndVersion,
		Name:    client.GetNoteName(noteID, content),
	}
	return client.writeNote(noteID, &note, content)
}

// given a deletion change, deletes the specified note
func (client *Client) applyDeletionChange(change *NoteChange) error {
	noteID := NoteID(change.EntityID)
	note, ok := client.storage.Notes[noteID]
	if !ok {
		return errors.New("note with id " + string(noteID) + " does not exist")
	}

	// remove file
	filename := client.getFileName(note.Name)
	vFilename := client.getVersionFileName(note.Name)
	// silently ignore errors
	os.Remove(filename)
	os.Remove(vFilename)

	// remove from cache
	delete(client.storage.Notes, noteID)

	return nil
}
