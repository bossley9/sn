package client

import (
	"encoding/json"
	"os"
	"path/filepath"

	f "git.sr.ht/~bossley9/sn/pkg/fileio"
)

type localStorage struct {
	filename       string
	filenameCompat string
	ChangeVersion  string `json:"cv"`
	AuthToken      string `json:"at"`
}

func newLocalStorage(name string) (*localStorage, error) {
	filename := getLocalStorageFilename(name)

	if err := os.MkdirAll(filepath.Dir(filename), f.RWX); err != nil {
		return nil, err
	}

	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			return nil, err
		}
		file.Close()
	} else if err != nil {
		return nil, err
	}

	file, _ := os.ReadFile(filename)

	storage := localStorage{}
	if len(file) > 0 {
		if err := json.Unmarshal(file, &storage); err != nil {
			return nil, err
		}
	}

	storage.filename = filename
	storage.filenameCompat = filename + "2"

	return &storage, nil
}

func getLocalStorageFilename(name string) string {
	storageDir := os.Getenv("XDG_DATA_HOME")
	if len(storageDir) == 0 {
		homeDir := os.Getenv("HOME")
		if len(homeDir) > 0 {
			storageDir = homeDir + "/.local/share"
		} else {
			storageDir = "~/.local/share"
		}
	}
	return storageDir + "/" + name + ".json"
}

func (storage *localStorage) writeChanges() error {
	raw, err := json.Marshal(storage)
	if err != nil {
		return err
	}
	return os.WriteFile(storage.filename, raw, f.RW)
}
