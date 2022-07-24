package client

import (
	"fmt"
	"os"

	j "git.sr.ht/~bossley9/sn/pkg/jsondiff"
	s "git.sr.ht/~bossley9/sn/pkg/simperium"
)

// given any change, applies that change to the specified note
func (client *client) applyChange(change *s.Change[NoteDiff]) {
	noteID := change.EntityID

	if change.Values.CreationDate.Operation == j.OP_INSERT {
		// note creation
		fmt.Println("\t\tcreating note " + noteID + "...")
		if err := client.applyCreationChange(change); err != nil {
			fmt.Println(err)
			fmt.Println("\t\tunable to create note " + noteID + ". Skipping...")
		}
	} else if change.Values.Deleted.Value {
		// note deletion
		fmt.Println("\t\tdeleting note " + noteID + "...")
		if err := client.applyDeletionChange(change); err != nil {
			fmt.Println(err)
			fmt.Println("\t\tunable to delete note " + noteID + ". Skipping...")
		}
	} else {
		// note update
		fmt.Println("\t\tupdating note " + noteID + "...")
		if err := client.applyUpdateChange(change); err != nil {
			fmt.Println(err)
			fmt.Println("\t\tunable to update note " + noteID + ". Skipping...")
		}
	}

	// update change version
	fmt.Println("\t\tupdating change version from " + client.cache.CurrentVersion + " to " + change.ChangeVersion + "...")
	if err := client.setCurrentVersion(change.ChangeVersion); err != nil {
		fmt.Println("\t\tunable to set current version. Skipping...")
		return
	}
}

// given an update change, applies that change to the specified note
func (client *client) applyUpdateChange(change *s.Change[NoteDiff]) error {
	noteID := change.EntityID

	content, err := client.readNote(noteID)
	if err != nil {
		return err
	}

	fmt.Println("\t\tapplying change " + change.ChangeVersion + " to note " + noteID + "...")
	result := change.Values.Content.Apply(string(content))

	fmt.Println("\t\twriting changes...")
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
func (client *client) applyCreationChange(change *s.Change[NoteDiff]) error {
	summary := NoteSummary{
		ID:      change.EntityID,
		Version: change.EndVersion,
		Content: change.Values.Content.Value,
	}
	return client.writeNote(&summary)
}

// given a deletion change, deletes the specified note
func (client *client) applyDeletionChange(change *s.Change[NoteDiff]) error {
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

	// remove from cache
	delete(client.cache.Notes, noteID)

	return nil
}
