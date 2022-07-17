package main

import (
	"fmt"
	"os"

	c "git.sr.ht/~bossley9/sn/pkg/client"
)

func main() {
	fmt.Println("initializing client...")
	_, err := c.NewClient()
	if err != nil {
		fmt.Println(err)
		fmt.Println("unable to initialize client. Exiting")
		os.Exit(1)
	}

	fmt.Println("done.")
}
