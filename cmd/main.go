package main

import (
	"fmt"
	"os"

	c "git.sr.ht/~bossley9/sn/pkg/client"
)

func main() {
	fmt.Println("initializing client...")
	client, err := c.NewClient()
	if err != nil {
		fmt.Println(err)
		fmt.Println("unable to initialize client. Exiting")
		os.Exit(1)
	}

	fmt.Println("authenticating with Simplenote...")
	if err := client.Authenticate(); err != nil {
		fmt.Println(err)
	}

	fmt.Println("done.")
}
