package localstorage

import (
	"encoding/json"
	"os"

	f "git.sr.ht/~bossley9/sn/pkg/fileio"
)

func (storage *LocalStorage) readFile() error {
	_, err := os.Stat(storage.filename)
	if os.IsNotExist(err) {
		file, err := os.Create(storage.filename)
		if err != nil {
			return err
		}
		file.Close()
	} else if err != nil {
		return err
	}

	raw, err := os.ReadFile(storage.filename)
	if err != nil {
		return err
	}

	var content map[string][]byte
	if raw == nil || len(raw) == 0 {
		content = map[string][]byte{}
	} else {
		if err := json.Unmarshal(raw, &content); err != nil {
			return err
		}
	}

	storage.content = content
	return nil
}

func (storage *LocalStorage) writeFile() error {
	byteContent, err := json.Marshal(storage.content)
	if err != nil {
		return err
	}
	return os.WriteFile(storage.filename, byteContent, f.RW)
}
