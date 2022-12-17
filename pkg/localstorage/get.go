package localstorage

import (
	"bytes"
	"encoding/gob"
	"errors"
)

func (storage *LocalStorage) Get(key string, value interface{}) error {
	storedValue := storage.content[key]
	if storedValue == nil {
		return errors.New("value for key '" + key + "' not found in storage.")
	}

	decoder := gob.NewDecoder(bytes.NewReader(storedValue))
	return decoder.Decode(value)
}
