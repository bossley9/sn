package client

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"

	"golang.org/x/term"
	"nhooyr.io/websocket"

	f "github.com/bossley9/sn/pkg/fileio"
	l "github.com/bossley9/sn/pkg/logger"
	s "github.com/bossley9/sn/pkg/simperium"
)

type Client struct {
	projectDir string
	versionDir string
	simp       *s.Client[NoteDiff]
	connection *websocket.Conn
	storage    *localStorage
}

func NewClient() (*Client, error) {
	c := Client{}

	storage, err := newLocalStorage("sn")
	if err != nil {
		return nil, err
	}
	c.storage = storage

	// initializing project directory
	home := os.Getenv("HOME")
	if len(home) == 0 {
		home = "."
	}
	c.projectDir = home + "/Documents/sn"
	if err := os.MkdirAll(c.projectDir, f.RWX); err != nil {
		return nil, err
	}

	// initializing version control
	// creating a directory within .git to automatically ignore version
	// metadata in most IDEs
	c.versionDir = c.projectDir + "/.git/version"
	if err := os.MkdirAll(c.versionDir, f.RWX); err != nil {
		return nil, err
	}

	// creating simperium client
	c.simp = s.NewClient[NoteDiff](APP_ID, API_KEY, "1.1", "node", "node-simperium", "0.0.1")

	return &c, nil
}

func fetchCredentials() (string, string, error) {
	l.PrintPlain("\n")
	l.PrintInfo("Username: ")
	var username string
	fmt.Scanln(&username)

	l.PrintInfo("Password (will not echo): ")
	password, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", "", err
	}
	l.PrintPlain("\n")
	return username, string(password), nil
}

// retrieve user authentication token
func (client *Client) Authenticate() error {
	if len(client.storage.AuthToken) > 0 {
		return nil
	}

	var username, password string
	var err error

	bwCmd := exec.Command("bw")
	if err := bwCmd.Run(); err != nil {
		l.PrintInfo("Bitwarden detected.\n")
		username, password, err = fetchBitwardenCredentials()
		if err != nil {
			username, password, err = fetchCredentials()
			if err != nil {
				return err
			}
		}
	} else {
		username, password, err = fetchCredentials()
		if err != nil {
			return err
		}
	}

	l.PrintInfo("Fetching authorization... ")
	token, err := client.simp.Authorize(username, password)
	if err != nil {
		return err
	}

	client.storage.AuthToken = token
	return nil
}

// connect to the server websocket
func (client *Client) Connect(ctx context.Context) error {
	return client.simp.ConnectToSocket(ctx)
}

// disconnect from the server websocket
func (client *Client) Disconnect() error {
	return client.simp.DisconnectSocket()
}

// authorize access to a given bucket
func (client *Client) OpenBucket(bucketName string, ctx context.Context) error {
	timedContext, cancel := context.WithTimeout(ctx, time.Second*4)
	defer cancel()
	errChan := make(chan error)

	go func() {
		authToken := client.storage.AuthToken
		if err := client.simp.WriteInitMessage(timedContext, 0, authToken, bucketName); err != nil {
			errChan <- err
		}

		// NOTE: This isn't in the Simperium documentation
		// the server sends two messages on initial auth
		if _, err := client.simp.ReadMessage(timedContext); err != nil {
			errChan <- err
		}
		if _, err := client.simp.ReadMessage(timedContext); err != nil {
			errChan <- err
		}

		errChan <- nil
	}()

	for {
		select {
		case <-timedContext.Done():
			return errors.New("bucket authorization timeout.")
		case err := <-errChan:
			return err
		}
	}
}
