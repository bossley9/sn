package main

import (
	"context"
	"os"

	sn "git.sr.ht/~bossley9/sn/pkg/sn"
)

func main() {
	args := os.Args
	arg := ""
	if len(args) > 1 {
		arg = args[1]
	}

	ctx := context.Background()

	sn.Run(arg, ctx)
}
