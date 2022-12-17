package localstorage

import (
	"os"
)

type localStorage struct {
	// TODO: make private when public read and write methods are created
	Filename string
	content  map[string][]byte
}

func New(name string) *localStorage {
	filename := getStorageFilename(name)

	return &localStorage{
		Filename: filename,
		content:  map[string][]byte{},
	}
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
