package client

import (
	"encoding/json"
	"errors"
	"os"

	s "git.sr.ht/~bossley9/sn/pkg/simperium"
)

type Cache struct {
	AuthToken      string               `json:"token"`
	CurrentVersion string               `json:"current_version"`
	Notes          map[string]NoteCache `json:"notes"`
}

type NoteCache struct {
	Version int    `json:"v"`
	Name    string `json:"n"`
}

func getCacheFile() string {
	cacheDir := os.Getenv("XDG_CACHE_HOME")
	if len(cacheDir) == 0 {
		cacheDir = "."
	}
	return cacheDir + "/snrc.json"
}

func ReadCache() (*Cache, error) {
	cacheFile := getCacheFile()

	file, err := os.ReadFile(cacheFile)
	if err != nil {
		return nil, err
	}

	var cache Cache
	if err := json.Unmarshal(file, &cache); err != nil {
		return nil, err
	}

	return &cache, nil
}

func (client *client) writeCache() error {
	cacheFile := getCacheFile()

	cacheContent, err := json.Marshal(client.cache)
	if err != nil {
		return err
	}
	if err := os.WriteFile(cacheFile, cacheContent, 0600); err != nil {
		return err
	}

	return nil
}

func (client *client) setCurrentVersion(version string) error {
	client.cache.CurrentVersion = version
	return client.writeCache()
}

func (client *client) saveNote(note *s.EntitySummary[Note]) error {
	if client.cache.Notes == nil {
		client.cache.Notes = make(map[string]NoteCache)
	}

	client.cache.Notes[note.ID] = NoteCache{
		Version: note.Version,
		Name:    GetNoteName(note.ID, &note.Data),
	}

	return client.writeCache()
}

func (client *client) setNoteVersion(noteID string, version int) error {
	if client.cache.Notes == nil {
		return errors.New("note cache does not exist")
	}
	note, ok := client.cache.Notes[noteID]
	if !ok {
		return errors.New("note with id " + noteID + " does not exist")
	}

	note.Version = version
	client.cache.Notes[noteID] = note

	return client.writeCache()
}

func (client *client) getCachedNote(noteID string) (NoteCache, error) {
	if client.cache.Notes == nil {
		return NoteCache{}, errors.New("note cache does not exist")
	}
	note, ok := client.cache.Notes[noteID]
	if !ok {
		return NoteCache{}, errors.New("note with id " + noteID + " does not exist")
	}

	return note, nil
}
