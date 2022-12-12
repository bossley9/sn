package client

import (
	"fmt"
	"os"
	"syscall"

	"github.com/gorilla/websocket"
	"golang.org/x/term"

	f "git.sr.ht/~bossley9/sn/pkg/fileio"
	s "git.sr.ht/~bossley9/sn/pkg/simperium"
)

type Client struct {
	projectDir string
	versionDir string
	cache      *Cache
	simp       *s.Client
	connection *websocket.Conn
}

func NewClient() (*Client, error) {
	c := Client{}

	// initializing project directory
	home := os.Getenv("HOME")
	if len(home) == 0 {
		home = "."
	}
	c.projectDir = home + "/Documents/sn"
	if err := os.MkdirAll(c.projectDir, f.RWX); err != nil {
		return nil, err
	}

	// reading cache
	cache, err := ReadCache()
	if err != nil {
		fmt.Println("\tunable to parse cache. Continuing...")
		cache = &Cache{}
	}
	c.cache = cache

	// initializing version control
	// creating a directory within .git to automatically ignore version
	// metadata in most IDEs
	c.versionDir = c.projectDir + "/.git/version"
	if err := os.MkdirAll(c.versionDir, f.RWX); err != nil {
		return nil, err
	}

	// creating simperium client
	c.simp = s.NewClient(APP_ID, API_KEY, "1.1", "node", "node-simperium", "0.0.1")

	return &c, nil
}

// retrieve user authentication token
func (client *Client) Authenticate() error {
	if len(client.getToken()) > 0 {
		// found cached token
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

	return client.setToken(token)
}

// connect to the server websocket
func (client *Client) Connect() error {
	return client.simp.ConnectToSocket()
}

// disconnect from the server websocket
func (client *Client) Disconnect() error {
	return client.simp.DisconnectSocket()
}

// authorize access to a given bucket
func (client *Client) OpenBucket(bucketName string) error {
	if err := client.simp.WriteInitMessage(0, client.cache.AuthToken, bucketName); err != nil {
		return err
	}

	// need to read two messages for some reason -
	// this isn't in the Simperium documentation
	if _, err := client.simp.ReadMessage(); err != nil {
		return err
	}
	if _, err := client.simp.ReadMessage(); err != nil {
		return err
	}

	return nil
}
