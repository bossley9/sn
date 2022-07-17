package client

import (
	"fmt"
	"os"
)

type client struct {
	projectDir string
	cache      *Cache
}

func NewClient() (*client, error) {
	c := client{}

	fmt.Println("\tinitializing project directory...")
	home := os.Getenv("HOME")
	if len(home) == 0 {
		home = "."
	}
	c.projectDir = home + "/Documents/simplenote"
	if err := os.MkdirAll(c.projectDir, 0700); err != nil {
		return nil, err
	}

	fmt.Println("\treading cache...")
	cache, err := ReadCache()
	if err != nil {
		fmt.Println("\tunable to parse cache. Continuing...")
		cache = &Cache{}
	}
	c.cache = cache

	return &c, nil
}
