package localstorage

func (storage *localStorage) Delete(key string) {
	storage.content[key] = nil
}
