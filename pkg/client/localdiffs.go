package client

import (
	"os"
	"strings"
	"time"

	"github.com/google/uuid"

	j "github.com/bossley9/sn/pkg/jsondiff"
	l "github.com/bossley9/sn/pkg/logger"
)

func (client *Client) GetLocalDiffs() []NoteChange {
	diffs := []NoteChange{}
	notes := client.storage.Notes

	if len(notes) == 0 {
		return diffs
	}

	filenames := map[string]string{}

	dirEntries, err := os.ReadDir(client.projectDir)
	if err != nil {
		l.PrintWarning("unable to read directory " + client.projectDir + ". Continuing...\n")
	} else {
		for _, entry := range dirEntries {
			if entry.IsDir() {
				continue
			}
			noteName := strings.TrimSuffix(entry.Name(), ".md")
			filenames[client.getFileName(noteName)] = noteName
		}
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

		delete(filenames, filename)

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

	for filename, noteName := range filenames {
		raw, err := os.ReadFile(filename)
		if err != nil {
			l.PrintWarning("unable to read file " + filename + ". Skipping...\n")
			continue
		}

		l.PrintInfo("\n" + noteName + " was created.")

		noteID := uuid.New().String()
		// ensure id uniqueness
		_, keyExists := client.storage.Notes[NoteID(noteID)]
		for keyExists {
			noteID = uuid.New().String()
			_, keyExists = client.storage.Notes[NoteID(noteID)]
		}

		note := Note{
			Name:    noteName,
			Version: 0,
		}
		client.storage.Notes[NoteID(noteID)] = note

		date := float32(time.Now().Unix())

		diffs = append(diffs, NoteChange{
			EntityID:  noteID,
			Operation: j.OP_MODIFY,
			Values: NoteDiff{
				CreationDate: j.Float32JSONDiff{
					Operation: j.OP_INSERT,
					Value:     date,
				},
				ModificationDate: j.Float32JSONDiff{
					Operation: j.OP_INSERT,
					Value:     date,
				},
				Content: j.StringJSONDiff{
					Operation: j.OP_INSERT,
					Value:     string(raw),
				},
				Tags: j.JSONDiff[[]string]{
					Operation: j.OP_INSERT,
					Value:     []string{},
				},
				SystemTags: j.JSONDiff[[]string]{
					Operation: j.OP_INSERT,
					Value:     []string{"markdown"},
				},
				Deleted: j.BoolJSONDiff{
					Operation: j.OP_INSERT,
					Value:     false,
				},
				ShareURL: j.StringJSONDiff{
					Operation: j.OP_INSERT,
					Value:     "",
				},
				PublishURL: j.StringJSONDiff{
					Operation: j.OP_INSERT,
					Value:     "",
				},
			},
		})
	}

	if len(diffs) > 0 {
		l.PrintPlain("\n")
	}

	return diffs
}
