package client

import (
	"os"

	j "git.sr.ht/~bossley9/sn/pkg/jsondiff"
	l "git.sr.ht/~bossley9/sn/pkg/logger"
)

func (client *Client) GetLocalDiffs() map[string]j.StringJSONDiff {
	diffs := make(map[string]j.StringJSONDiff, 0)
	notes := client.cache.Notes

	if len(notes) == 0 {
		return diffs
	}

	for noteID, noteCache := range notes {
		filename := client.getFileName(noteCache.Name)
		vFilename := client.getVersionFileName(noteCache.Name)
		s1, err := os.ReadFile(vFilename)
		if err != nil {
			l.PrintWarning("unable to read file " + vFilename + " for versioning. Skipping...\n")
			continue
		}
		s2, err := os.ReadFile(filename)
		if err != nil {
			l.PrintWarning("unable to read file " + filename + " for versioning. Skipping...\n")
			continue
		}

		diff := j.GetDiff(string(s1), string(s2))
		if len(diff.Value) == 0 {
			continue // no diff found
		}

		l.PrintInfo("\nLocal diff found for " + noteCache.Name + ".")

		diffs[noteID] = diff
	}

	if len(diffs) > 0 {
		l.PrintPlain("\n")
	}

	return diffs
}
