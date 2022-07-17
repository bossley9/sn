package client

import (
	"fmt"
	"os"
	"syscall"

	"golang.org/x/term"

	s "git.sr.ht/~bossley9/sn/pkg/simperium"
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

func (client *client) Authenticate() error {
	if len(client.cache.AuthToken) > 0 {
		return nil
	}

	fmt.Print("\tusername: ")
	var username string
	fmt.Scanln(&username)

	fmt.Print("\tpassword (will not echo): ")
	password, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return err
	}
	fmt.Println()

	token, err := s.Authorize(APP_ID, API_KEY, username, string(password))
	if err != nil {
		return err
	}

	client.cache.AuthToken = token

	if err := WriteCache(client.cache); err != nil {
		return err
	}

	return nil
}
