package client

import (
	"os"

	l "git.sr.ht/~bossley9/sn/pkg/logger"
)

// clear and remove all data
func (client *Client) Clear() error {
	l.PrintInfo("Deleting all notes... ")
	if err := os.RemoveAll(client.projectDir); err != nil {
		return err
	}

	// delete cache file
	l.PrintInfo("deleting cache... ")
	// ignore error if file does not exist
	os.Remove(client.storage.filenameCompat)
	os.Remove(client.storage.filename)

	return nil
}
