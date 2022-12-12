// vim: fdm=indent
package sn

import (
	"os"

	c "git.sr.ht/~bossley9/sn/pkg/client"
	l "git.sr.ht/~bossley9/sn/pkg/logger"
)

func printFatalAndExit(err error) {
	l.PrintError(err)
	l.PrintError("\nFatal error. Exiting.\n")
	os.Exit(1)
}

func OpenProjectDir() {
	l.PrintInfo("Initializing client... ")
	client, err := c.NewClient()
	if err != nil {
		printFatalAndExit(err)
	}
	l.PrintInfo("done.\n")

	l.PrintInfo("Authenticating with server... ")
	if err := client.Authenticate(); err != nil {
		printFatalAndExit(err)
	}
	l.PrintInfo("done.\n")

	l.PrintInfo("Connecting to socket... ")
	if err := client.Connect(); err != nil {
		printFatalAndExit(err)
	}
	l.PrintInfo("done.\n")

	l.PrintInfo("Accessing notes... ")
	if err := client.OpenBucket("note"); err != nil {
		printFatalAndExit(err)
	}
	l.PrintInfo("done.\n")

	l.PrintInfo("Syncing client... ")
	if err := client.Sync(); err != nil {
		printFatalAndExit(err)
	}
	client.Disconnect() // disconnect after sync to prevent timeout
	l.PrintInfo("done.\n")

	// open project
	if err := client.OpenProjectDir(); err != nil {
		printFatalAndExit(err)
	}

	l.PrintInfo("Searching for local diffs... ")
	diffs := client.GetLocalDiffs()
	if len(diffs) == 0 {
		l.PrintWarning("no local diffs found.")
		l.PrintPlain("\n")
		os.Exit(0)
	}
	l.PrintInfo("done.\n")

	// reconnect after edits
	l.PrintInfo("Authenticating with server... ")
	if err := client.Authenticate(); err != nil {
		printFatalAndExit(err)
	}
	l.PrintInfo("done.\n")

	l.PrintInfo("Connecting to socket... ")
	if err := client.Connect(); err != nil {
		printFatalAndExit(err)
	}
	defer client.Disconnect()
	l.PrintInfo("done.\n")

	l.PrintInfo("Accessing notes... ")
	if err := client.OpenBucket("note"); err != nil {
		printFatalAndExit(err)
	}
	l.PrintInfo("done.\n")

	l.PrintInfo("Uploading diffs... ")
	if err := client.Upload(diffs); err != nil {
		printFatalAndExit(err)
	}
	l.PrintInfo("done.\n")
}

func Clear() {
	l.PrintInfo("Initializing client... ")
	client, err := c.NewClient()
	if err != nil {
		printFatalAndExit(err)
	}
	l.PrintInfo("done.\n")

	l.PrintInfo("Clearing data... ")
	if err := client.Clear(); err != nil {
		printFatalAndExit(err)
	}
	l.PrintInfo("done.\n")
}

func DownloadSync(reset bool) {
	l.PrintInfo("Initializing client... ")
	client, err := c.NewClient()
	if err != nil {
		printFatalAndExit(err)
	}
	l.PrintInfo("done.\n")

	l.PrintInfo("Authenticating with server... ")
	if err := client.Authenticate(); err != nil {
		printFatalAndExit(err)
	}
	l.PrintInfo("done.\n")

	l.PrintInfo("Connecting to socket... ")
	if err := client.Connect(); err != nil {
		printFatalAndExit(err)
	}
	defer client.Disconnect()
	l.PrintInfo("done.\n")

	l.PrintInfo("Accessing notes... ")
	if err := client.OpenBucket("note"); err != nil {
		printFatalAndExit(err)
	}
	l.PrintInfo("done.\n")

	if reset {
		l.PrintInfo("Refetching... ")
		if err := client.RefetchSync(); err != nil {
			printFatalAndExit(err)
		}
	} else {
		l.PrintInfo("Syncing client... ")
		if err := client.Sync(); err != nil {
			printFatalAndExit(err)
		}
	}
	l.PrintInfo("done.\n")
}

func UploadSync() {
	l.PrintInfo("Initializing client... ")
	client, err := c.NewClient()
	if err != nil {
		printFatalAndExit(err)
	}
	l.PrintInfo("done.\n")

	l.PrintInfo("Searching for local diffs... ")
	diffs := client.GetLocalDiffs()
	if len(diffs) == 0 {
		l.PrintWarning("no local diffs found.")
		l.PrintPlain("\n")
		os.Exit(0)
	}
	l.PrintInfo("done.\n")

	l.PrintInfo("Authenticating with server... ")
	if err := client.Authenticate(); err != nil {
		printFatalAndExit(err)
	}
	l.PrintInfo("done.\n")

	l.PrintInfo("Connecting to socket... ")
	if err := client.Connect(); err != nil {
		printFatalAndExit(err)
	}
	defer client.Disconnect()
	l.PrintInfo("done.\n")

	l.PrintInfo("Accessing notes... ")
	if err := client.OpenBucket("note"); err != nil {
		printFatalAndExit(err)
	}
	l.PrintInfo("done.\n")

	l.PrintInfo("Uploading diffs... ")
	if err := client.Upload(diffs); err != nil {
		printFatalAndExit(err)
	}
	l.PrintInfo("done.\n")
}
