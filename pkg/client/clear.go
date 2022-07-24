package client

import (
	"fmt"
	"os"
)

// clear and remove all data
func (client *client) Clear() error {
	fmt.Println("\tdeleting all notes...")
	if err := os.RemoveAll(client.projectDir); err != nil {
		return err
	}

	// delete cache file
	fmt.Println("\tdeleting cache...")
	if err := os.Remove(getCacheFile()); err != nil {
		return err
	}

	return nil
}
