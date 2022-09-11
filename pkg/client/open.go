package client

import (
	"os"
	"os/exec"
)

func (client *client) OpenProjectDir() error {
	cmd := exec.Command(os.Getenv("EDITOR"))
	cmd.Dir = client.projectDir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
