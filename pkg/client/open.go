package client

import (
	"os"
	"os/exec"
)

func (client *Client) OpenProjectDir() error {
	cmd := exec.Command(os.Getenv("EDITOR"))
	cmd.Dir = client.projectDir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
