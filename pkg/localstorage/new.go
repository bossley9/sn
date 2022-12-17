package localstorage

import (
	"os"
	"path/filepath"

	f "git.sr.ht/~bossley9/sn/pkg/fileio"
)

type LocalStorage struct {
	// TODO: make private when public read and write methods are created
	Filename string
	content  map[string][]byte
}

func New(name string) (*LocalStorage, error) {
	filename := getStorageFilename(name)

	if err := os.MkdirAll(filepath.Dir(filename), f.RWX); err != nil {
		return nil, err
	}

	return &LocalStorage{
		Filename: filename,
		content:  map[string][]byte{},
	}, nil
}

func getStorageFilename(name string) string {
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
