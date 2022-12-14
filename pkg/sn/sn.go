package sn

import (
	"context"
	"os"

	c "git.sr.ht/~bossley9/sn/pkg/client"
	l "git.sr.ht/~bossley9/sn/pkg/logger"
)

func Run(arg string, ctx context.Context) {
	l.PrintPlain("\n")
	var err error

	switch arg {
	case "":
		err = openProjectDir(ctx)
	case "c":
		err = clear()
	case "d":
		err = downloadSync(ctx)
	case "h":
		printUsage()
	case "r":
		err = refetchSync(ctx)
	case "u":
		err = uploadSync(ctx)
	default:
		printUsage()
	}

	if err != nil {
		l.PrintError(err)
		l.PrintError("\nFatal error. Exiting.\n")
		os.Exit(1)
	}

	l.PrintPlain("\n") // ensure terminal fg has reverted back to default color
}

func printUsage() {
	args := []string{
		"[no arg]  download and sync with server, then open the project directory with $EDITOR",
		"c         clear auth, reset cache and remove all notes",
		"d         download and sync with server",
		"h         view this help usage",
		"r         refetch all notes",
		"u         upload and sync with server",
	}

	l.PrintInfo("Usage: sn [arg]\n")
	for _, arg := range args {
		l.PrintInfo(arg + "\n")
	}
}

func initializeClient() (*c.Client, error) {
	l.PrintInfo("Initializing client... ")
	client, err := c.NewClient()
	if err != nil {
		return nil, err
	}
	l.PrintInfo("done.\n")

	return client, nil
}

func authenticateAndConnect(client *c.Client, ctx context.Context) error {
	l.PrintInfo("Authenticating with server... ")
	if err := client.Authenticate(); err != nil {
		return err
	}
	l.PrintInfo("done.\n")

	l.PrintInfo("Connecting to socket... ")
	if err := client.Connect(); err != nil {
		return err
	}
	l.PrintInfo("done.\n")

	l.PrintInfo("Accessing notes... ")
	const retryLimit = 3
	for i := 1; i <= retryLimit; i++ {
		err := client.OpenBucket("note", ctx)
		if err == nil {
			break
		}
		if i == retryLimit {
			return err
		}

		l.PrintWarning(err)
		l.PrintInfo("\nRetrying... ")
	}
	l.PrintInfo("done.\n")

	return nil
}

func uploadAvailableDiffs(client *c.Client, ctx context.Context) error {
	l.PrintInfo("Searching for local diffs... ")
	diffs := client.GetLocalDiffs()
	if len(diffs) == 0 {
		l.PrintWarning("no local diffs found.\n")
		return nil
	}
	l.PrintInfo("done.\n")

	if err := authenticateAndConnect(client, ctx); err != nil {
		return err
	}
	defer client.Disconnect()

	l.PrintInfo("Uploading diffs... ")
	if err := client.Upload(diffs); err != nil {
		return err
	}
	l.PrintInfo("done.\n")

	return nil
}

func openProjectDir(ctx context.Context) error {
	client, err := initializeClient()
	if err != nil {
		return err
	}

	if err := authenticateAndConnect(client, ctx); err != nil {
		return err
	}

	l.PrintInfo("Syncing client... ")
	if err := client.Sync(); err != nil {
		return err
	}
	l.PrintInfo("done.\n")

	client.Disconnect() // disconnect after sync to prevent timeout

	// open project
	if err := client.OpenProjectDir(); err != nil {
		return err
	}

	return uploadAvailableDiffs(client, ctx)
}

func clear() error {
	client, err := initializeClient()
	if err != nil {
		return err
	}

	l.PrintInfo("Clearing data... ")
	if err := client.Clear(); err != nil {
		return err
	}
	l.PrintInfo("done.\n")
	return nil
}

func downloadSync(ctx context.Context) error {
	client, err := initializeClient()
	if err != nil {
		return err
	}

	if err := authenticateAndConnect(client, ctx); err != nil {
		return err
	}
	defer client.Disconnect()

	l.PrintInfo("Syncing client... ")
	if err := client.Sync(); err != nil {
		return err
	}
	l.PrintInfo("done.\n")

	return nil
}

func refetchSync(ctx context.Context) error {
	client, err := initializeClient()
	if err != nil {
		return err
	}

	if err := authenticateAndConnect(client, ctx); err != nil {
		return err
	}
	defer client.Disconnect()

	l.PrintInfo("Refetching... ")
	if err := client.RefetchSync(); err != nil {
		return err
	}
	l.PrintInfo("done.\n")

	return nil
}

func uploadSync(ctx context.Context) error {
	client, err := initializeClient()
	if err != nil {
		return err
	}
	return uploadAvailableDiffs(client, ctx)
}
