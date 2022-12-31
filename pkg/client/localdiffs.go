package client

import (
	"os"
	"time"

	j "git.sr.ht/~bossley9/sn/pkg/jsondiff"
	l "git.sr.ht/~bossley9/sn/pkg/logger"
)

func (client *Client) GetLocalDiffs() []NoteChange {
	diffs := []NoteChange{}
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
			l.PrintInfo("\n" + note.Name + " was deleted.")

			// NOTE: This isn't in the Simperium documentation
			// deleting notes requires two diffs: one for moving the
			// note to trash and a second for deleting the note
			diffs = append(diffs, NoteChange{
				EntityID:  string(noteID),
				Operation: j.OP_MODIFY,
				Values: NoteDiff{
					Deleted: j.BoolJSONDiff{
						Operation: j.OP_REPLACE,
						Value:     true,
					},
					ModificationDate: j.Float32JSONDiff{
						Operation: j.OP_REPLACE,
						Value:     float32(time.Now().Unix()),
					},
				},
			})

			diffs = append(diffs, NoteChange{
				EntityID:  string(noteID),
				Operation: j.OP_DELETE,
			})

			continue
		}

		contentDiff := j.GetDiff(string(s1), string(s2))
		if len(contentDiff.Value) == 0 {
			continue
		}

		l.PrintInfo("\nLocal diff found for " + note.Name + ".")

		diffs = append(diffs, NoteChange{
			EntityID:  string(noteID),
			Operation: j.OP_MODIFY,
			Values: NoteDiff{
				Content: contentDiff,
				ModificationDate: j.Float32JSONDiff{
					Operation: j.OP_REPLACE,
					Value:     float32(time.Now().Unix()),
				},
			},
		})
	}

	if len(diffs) > 0 {
		l.PrintPlain("\n")
	}

	return diffs
}
