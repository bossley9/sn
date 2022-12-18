package client

import "errors"

func (client *Client) setNoteVersion(noteID string, version int) error {
	if client.storage.Notes == nil {
		return errors.New("note cache does not exist")
	}
	note, ok := client.storage.Notes[noteID]
	if !ok {
		return errors.New("note with id " + noteID + " does not exist")
	}

	note.Version = version
	client.storage.Notes[noteID] = note

	return client.storage.writeChanges()
}

func (client *Client) getCachedNote(noteID string) (NoteCache, error) {
	if client.storage.Notes == nil {
		return NoteCache{}, errors.New("note cache does not exist")
	}
	note, ok := client.storage.Notes[noteID]
	if !ok {
		return NoteCache{}, errors.New("note with id " + noteID + " does not exist")
	}

	return note, nil
}
