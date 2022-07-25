package client

import (
	"os"
	"os/exec"
)

func (client *client) OpenDirectory() error {
	cmd := exec.Command("/usr/local/bin/nvim", "-R")
	cmd.Dir = client.projectDir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
