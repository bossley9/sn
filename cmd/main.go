// vim: fdm=indent
package main

import (
	"fmt"
	"os"

	sn "git.sr.ht/~bossley9/sn/pkg/sn"
)

const cyan = "\033[0;36m"
const none = "\033[0m"

func printusage() {
	args := []string{
		"[no arg]  download and sync with server, then open the project directory with $EDITOR",
		"c         clear auth, reset cache and remove all notes",
		"d         download and sync with server",
		"h         view this help usage",
		"r         refetch all notes",
		"u         upload and sync with server",
	}

	fmt.Println(cyan)
	fmt.Println("Usage: sn [arg]")
	for _, arg := range args {
		fmt.Println(arg)
	}
	fmt.Println(none)
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
		sn.Downloadsync(false)
	case "h":
		printusage()
	case "r":
		sn.Downloadsync(true)
	case "u":
		sn.Uploadsync()
	default:
		printusage()
		return
	}
}
