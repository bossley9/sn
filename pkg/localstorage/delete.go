package localstorage

func (storage *LocalStorage) Delete(key string) error {
	storage.content[key] = nil
	return storage.writeFile()
}
