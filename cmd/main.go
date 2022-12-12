// vim: fdm=indent
package main

import (
	"os"

	l "git.sr.ht/~bossley9/sn/pkg/logger"
	sn "git.sr.ht/~bossley9/sn/pkg/sn"
)

func printUsage() {
	args := []string{
		"[no arg]  download and sync with server, then open the project directory with $EDITOR",
		"c         clear auth, reset cache and remove all notes",
		"d         download and sync with server",
		"h         view this help usage",
		"r         refetch all notes",
		"u         upload and sync with server",
	}

	l.PrintInfo("\n")
	l.PrintInfo("Usage: sn [arg]\n")
	for _, arg := range args {
		l.PrintInfo(arg + "\n")
	}
	l.PrintInfo("\n")
}

func main() {
	args := os.Args
	arg := ""
	if len(args) > 1 {
		arg = args[1]
	}

	switch arg {
	case "":
		sn.OpenProjectDir()
	case "c":
		sn.Clear()
	case "d":
		sn.DownloadSync(false)
	case "h":
		printUsage()
	case "r":
		sn.DownloadSync(true)
	case "u":
		sn.UploadSync()
	default:
		printUsage()
	}
}
