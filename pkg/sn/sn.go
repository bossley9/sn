// vim: fdm=indent
package sn

import (
	"fmt"
	"log"

	c "git.sr.ht/~bossley9/sn/pkg/client"
)

const red = "\033[0;31m"
const cyan = "\033[0;36m"
const none = "\033[0m"

func printFatalAndExit(err error) {
	fmt.Print(red)
	fmt.Println(err)
	log.Fatal("Fatal error. Exiting." + none)
}

func OpenProjectDir() {
	fmt.Println("initializing client...")
	client, err := c.NewClient()
	if err != nil {
		fmt.Println(err)
		log.Fatal("unable to initialize client. Exiting.")
	}

	fmt.Println("authenticating with server...")
	if err := client.Authenticate(); err != nil {
		fmt.Println(err)
	}

	fmt.Println("connecting to socket...")
	if err := client.Connect(); err != nil {
		fmt.Println(err)
	}

	fmt.Println("accessing notes...")
	if err := client.OpenBucket("note"); err != nil {
		fmt.Println(err)
	}

	fmt.Println("syncing client...")
	if err := client.Sync(); err != nil {
		fmt.Println(err)
	}

	client.Disconnect() // disconnect after sync to prevent timeout

	// open project
	if err := client.OpenProjectDir(); err != nil {
		fmt.Println(err)
		log.Fatal("unable to open $EDITOR. Exiting.")
	}

	fmt.Println("searching for local diffs...")
	diffs, err := client.GetLocalDiffs()
	if err != nil {
		fmt.Println(err)
		log.Fatal("unable to find local diffs. Exiting.")
	}
	if len(diffs) == 0 {
		fmt.Println("no local diffs. Exiting.")
		return
	}

	// reconnect after edits
	fmt.Println("authenticating with server...")
	if err := client.Authenticate(); err != nil {
		fmt.Println(err)
	}

	fmt.Println("connecting to socket...")
	if err := client.Connect(); err != nil {
		fmt.Println(err)
	}

	defer client.Disconnect()

	fmt.Println("accessing notes...")
	if err := client.OpenBucket("note"); err != nil {
		fmt.Println(err)
	}

	fmt.Println("uploading diffs...")
	if err := client.Upload(diffs); err != nil {
		fmt.Println(err)
		log.Fatal("unable to upload diffs. Exiting.")
	}
}

func Clear() {
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

func DownloadSync(reset bool) {
	fmt.Println(cyan)

	fmt.Print("Initializing client... ")
	client, err := c.NewClient()
	if err != nil {
		printFatalAndExit(err)
	}
	fmt.Println("done.")

	fmt.Print("Authenticating with server... ")
	if err := client.Authenticate(); err != nil {
		printFatalAndExit(err)
	}
	fmt.Println("done.")

	fmt.Print("Connecting to socket... ")
	if err := client.Connect(); err != nil {
		printFatalAndExit(err)
	}
	defer client.Disconnect()
	fmt.Println("done.")

	fmt.Print("Accessing notes... ")
	if err := client.OpenBucket("note"); err != nil {
		printFatalAndExit(err)
	}
	fmt.Print(cyan)
	fmt.Println("done.")

	if reset {
		fmt.Print("Refetching...")
		if err := client.RefetchSync(); err != nil {
			printFatalAndExit(err)
		}
	} else {
		fmt.Print("Syncing client... ")
		if err := client.Sync(); err != nil {
			printFatalAndExit(err)
		}
	}
	fmt.Print(cyan)
	fmt.Println("done.")

	fmt.Println(none)
}

func UploadSync() {
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
