package localstorage

import (
	"bytes"
	"encoding/gob"
)

func (storage *LocalStorage) Set(key string, value interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(value); err != nil {
		return err
	}

	storage.content[key] = buf.Bytes()
	return nil
}
