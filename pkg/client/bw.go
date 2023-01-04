package client

import (
	"os/exec"
	"strings"
	"syscall"

	"golang.org/x/term"

	l "git.sr.ht/~bossley9/sn/pkg/logger"
)

func fetchBitwardenCredentials() (string, string, error) {
	l.PrintInfo("Bitwarden Master Password (will not echo): ")
	rawMasterPassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", "", err
	}
	l.PrintPlain("\n")

	const BW_SN_ID = "simplenote"

	usernameCmd := exec.Command("bw", "get", "username", BW_SN_ID)
	usernameCmd.Stdin = strings.NewReader(string(rawMasterPassword))
	rawUsername, err := usernameCmd.Output()
	if err != nil {
		return "", "", err
	}

	passwordCmd := exec.Command("bw", "get", "password", BW_SN_ID)
	passwordCmd.Stdin = strings.NewReader(string(rawMasterPassword))
	rawPassword, err := passwordCmd.Output()
	if err != nil {
		return "", "", err
	}

	return string(rawUsername), string(rawPassword), nil
}
