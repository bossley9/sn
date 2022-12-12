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
	if err := os.Remove(getCacheFile()); err != nil {
		return err
	}

	return nil
}
