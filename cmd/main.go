package main

import (
	"fmt"
	"log"
	"os"

	c "git.sr.ht/~bossley9/sn/pkg/client"
)

func printusage() {
	usage :=
		`usage: sn [d]
	[no arg]  same as using the argument "d"
	c         clear auth, reset cache and remove all notes
	d         download and sync with server
	r         reset cache and refetch all notes
	u         upload and sync with server`
	fmt.Println(usage)
}

func main() {
	args := os.Args
	arg := "d"

	if len(args) > 1 {
		arg = args[1]
	}

	switch arg {
	case "c":
		clear()
	case "d":
		downloadsync(false)
	case "r":
		downloadsync(true)
	case "u":
		uploadsync()
	default:
		printusage()
		return
	}

	fmt.Println("done.")
}

func clear() {
	fmt.Println("initializing client...")
	client, err := c.NewClient()
	if err != nil {
		fmt.Println(err)
		log.Fatal("unable to initialize client. Exiting.")
	}

	fmt.Println("clearing data...")
	if err := client.Clear(); err != nil {
		fmt.Println(err)
	}
}

func downloadsync(reset bool) {
	fmt.Println("initializing client...")
	client, err := c.NewClient()
	if err != nil {
		fmt.Println(err)
		log.Fatal("unable to initialize client. Exiting.")
	}

	fmt.Println("authenticating with server...")
	if err := client.Authenticate(); err != nil {
		fmt.Println(err)
		log.Fatal("unable to authenticate. Exiting.")
	}

	fmt.Println("connecting to socket...")
	if err := client.Connect(); err != nil {
		fmt.Println(err)
		log.Fatal("unable to connect to socket. Exiting.")
	}

	defer client.Disconnect()

	fmt.Println("accessing notes...")
	if err := client.OpenBucket("note"); err != nil {
		fmt.Println(err)
		log.Fatal("unable to open bucket. Exiting.")
	}

	if reset {
		fmt.Println("refetching...")
		if err := client.RefetchSync(); err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Println("syncing client...")
		if err := client.Sync(); err != nil {
			fmt.Println(err)
		}
	}
}

func uploadsync() {
	fmt.Println("initializing client...")
	client, err := c.NewClient()
	if err != nil {
		fmt.Println(err)
		log.Fatal("unable to initialize client. Exiting.")
	}

	fmt.Println("searching for local diffs...")
	diffs, err := client.GetLocalDiffs()
	if err != nil {
		fmt.Println(err)
		log.Fatal("unable to find local diffs. Exiting.")
	}
	if len(diffs) == 0 {
		log.Fatal("no local diffs. Exiting.")
	}

	fmt.Println("authenticating with server...")
	if err := client.Authenticate(); err != nil {
		fmt.Println(err)
		log.Fatal("unable to authenticate. Exiting.")
	}

	fmt.Println("connecting to socket...")
	if err := client.Connect(); err != nil {
		fmt.Println(err)
		log.Fatal("unable to connect to socket. Exiting.")
	}

	defer client.Disconnect()

	fmt.Println("accessing notes...")
	if err := client.OpenBucket("note"); err != nil {
		fmt.Println(err)
		log.Fatal("unable to open bucket. Exiting.")
	}

	fmt.Println("uploading diffs...")
	if err := client.Upload(diffs); err != nil {
		fmt.Println(err)
		log.Fatal("unable to upload diffs. Exiting.")
	}
}
