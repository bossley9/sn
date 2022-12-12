// vim: fdm=indent
package sn

import (
	"fmt"
	"log"

	c "git.sr.ht/~bossley9/sn/pkg/client"
)

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
