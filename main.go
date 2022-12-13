package main

import (
	"os"

	sn "git.sr.ht/~bossley9/sn/pkg/sn"
)

func main() {
	args := os.Args
	arg := ""
	if len(args) > 1 {
		arg = args[1]
	}

	sn.Run(arg)
}
