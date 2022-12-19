package client

import (
	"os"

	j "git.sr.ht/~bossley9/sn/pkg/jsondiff"
	l "git.sr.ht/~bossley9/sn/pkg/logger"
)

func (client *Client) GetLocalDiffs() map[NoteID]j.StringJSONDiff {
	diffs := make(map[NoteID]j.StringJSONDiff, 0)
	notes := client.storage.Notes

	if len(notes) == 0 {
		return diffs
	}

	for noteID, note := range notes {
		filename := client.getFileName(note.Name)
		vFilename := client.getVersionFileName(note.Name)
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

		l.PrintInfo("\nLocal diff found for " + note.Name + ".")

		diffs[noteID] = diff
	}

	if len(diffs) > 0 {
		l.PrintPlain("\n")
	}

	return diffs
}
