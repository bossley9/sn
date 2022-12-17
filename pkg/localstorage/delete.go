package localstorage

func (storage *LocalStorage) Delete(key string) {
	storage.content[key] = nil
}
