package client

import "errors"

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
