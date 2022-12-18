package localstorage

import (
	"os"
	"path/filepath"

	f "git.sr.ht/~bossley9/sn/pkg/fileio"
)

type LocalStorage struct {
	filename string
	content  map[string][]byte
	// TODO remove when migrated off cache utils
	FilenameCompat string
}

func New(name string) (*LocalStorage, error) {
	filename := getStorageFilename(name)

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

	storage := LocalStorage{
		FilenameCompat: filename + "2",
		filename:       filename,
		content:        map[string][]byte{},
	}

	if err := storage.readFile(); err != nil {
		return nil, err
	}

	return &storage, nil
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
