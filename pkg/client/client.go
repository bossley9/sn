package client

import (
	"fmt"
	"os"
	"syscall"

	"github.com/gorilla/websocket"
	"golang.org/x/term"

	s "git.sr.ht/~bossley9/sn/pkg/simperium"
)

type client struct {
	projectDir string
	cache      *Cache
	simp       *s.Client
	connection *websocket.Conn
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

	fmt.Println("\tcreating simperium client...")
	c.simp = s.NewClient(APP_ID, API_KEY, "1.1", "node", "node-simperium", "0.0.1")

	return &c, nil
}

// retrieve user authentication token
func (client *client) Authenticate() error {
	if len(client.cache.AuthToken) > 0 {
		fmt.Println("\tfound cached token.")
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

	fmt.Println("\tfetching authorization...")
	token, err := client.simp.Authorize(username, string(password))
	if err != nil {
		return err
	}

	client.cache.AuthToken = token

	if err := client.writeCache(); err != nil {
		return err
	}

	return nil
}

// connect to the server websocket
func (client *client) Connect() error {
	if err := client.simp.ConnectToSocket(); err != nil {
		return err
	}
	return nil
}

// disconnect from the server websocket
func (client *client) Disconnect() error {
	if err := client.simp.DisconnectSocket(); err != nil {
		return err
	}
	return nil
}

// authorize access to a given bucket
func (client *client) OpenBucket(bucketName string) error {
	if err := client.simp.WriteInitMessage(0, client.cache.AuthToken, bucketName); err != nil {
		return err
	}

	if _, err := client.simp.ReadMessage(); err != nil {
		return err
	}
	if _, err := client.simp.ReadMessage(); err != nil {
		return err
	}

	return nil
}
