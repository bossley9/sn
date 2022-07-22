package client

import (
	"encoding/json"
	"os"
)

type Cache struct {
	AuthToken      string               `json:"token"`
	CurrentVersion string               `json:"current_version"`
	Notes          map[string]NoteCache `json:"notes"`
}

type NoteCache struct {
	Version  int    `json:"v"`
	Filename string `json:"fn"`
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

func (client *client) saveNote(entityID string, version int, filename string) error {
	if client.cache.Notes == nil {
		client.cache.Notes = make(map[string]NoteCache)
	}

	client.cache.Notes[entityID] = NoteCache{
		Version:  version,
		Filename: filename,
	}

	return client.writeCache()
}
