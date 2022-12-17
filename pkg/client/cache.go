package client

import (
	"encoding/json"
	"errors"
	"os"

	f "git.sr.ht/~bossley9/sn/pkg/fileio"
	ls "git.sr.ht/~bossley9/sn/pkg/localstorage"
)

type Cache struct {
	AuthToken      string               `json:"t"`
	CurrentVersion string               `json:"cv"`
	Notes          map[string]NoteCache `json:"n"`
}

type NoteCache struct {
	Version int    `json:"v"`
	Name    string `json:"n"`
}

func getCacheFile() (string, error) {
	storage, err := ls.New("sn")
	if err != nil {
		return "", err
	}
	return storage.Filename, nil
}

func ReadCache() (*Cache, error) {
	cacheFile, err := getCacheFile()
	if err != nil {
		return nil, err
	}

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

func (client *Client) writeCache() error {
	cacheFile, err := getCacheFile()
	if err != nil {
		return err
	}

	cacheContent, err := json.Marshal(client.cache)
	if err != nil {
		return err
	}
	if err := os.WriteFile(cacheFile, cacheContent, f.RW); err != nil {
		return err
	}

	return nil
}

func (client *Client) getToken() string {
	return client.cache.AuthToken
}
func (client *Client) setToken(token string) error {
	client.cache.AuthToken = token
	return client.writeCache()
}

func (client *Client) getCurrentVersion() string {
	return client.cache.CurrentVersion
}
func (client *Client) setCurrentVersion(version string) error {
	client.cache.CurrentVersion = version
	return client.writeCache()
}

func (client *Client) setNoteVersion(noteID string, version int) error {
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

func (client *Client) getCachedNote(noteID string) (NoteCache, error) {
	if client.cache.Notes == nil {
		return NoteCache{}, errors.New("note cache does not exist")
	}
	note, ok := client.cache.Notes[noteID]
	if !ok {
		return NoteCache{}, errors.New("note with id " + noteID + " does not exist")
	}

	return note, nil
}
