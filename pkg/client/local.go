package client

import (
	"fmt"
	"os"

	j "git.sr.ht/~bossley9/sn/pkg/jsondiff"
)

func (client *Client) GetLocalDiffs() (map[string]j.StringJSONDiff, error) {
	diffs := make(map[string]j.StringJSONDiff, 0)
	notes := client.cache.Notes

	if len(notes) == 0 {
		return diffs, nil
	}

	for noteID, noteCache := range notes {
		filename := client.getFileName(noteCache.Name)
		vFilename := client.getVersionFileName(noteCache.Name)
		s1, err := os.ReadFile(vFilename)
		if err != nil {
			fmt.Println(err)
			fmt.Println("\t\tunable to read file " + vFilename + " for versioning. Skipping...")
			continue
		}
		s2, err := os.ReadFile(filename)
		if err != nil {
			fmt.Println(err)
			fmt.Println("\t\tunable to read file " + filename + " for versioning. Skipping...")
			continue
		}

		diff := j.GetDiff(string(s1), string(s2))
		if len(diff.Value) == 0 {
			continue // no diff found
		}

		fmt.Println("\t\tdiff found for " + noteCache.Name + ".")

		diffs[noteID] = diff
	}

	return diffs, nil
}
