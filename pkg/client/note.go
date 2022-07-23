package client

import (
	"os"
	"strings"

	j "git.sr.ht/~bossley9/sn/pkg/jsondiff"
	s "git.sr.ht/~bossley9/sn/pkg/simperium"
)

type Note struct {
	Tags           []string `json:"tags"`
	Deleted        bool     `json:"deleted"`
	ShareURL       string   `json:"shareURL"`
	PublishURL     string   `json:"publishURL"`
	Content        string   `json:"content"`
	SystemTags     []string `json:"systemTags"`
	LastEditedDate float32  `json:"modificationDate"`
	CreationDate   float32  `json:"creationDate"`
}

type NoteDiff struct {
	Content j.StringJSONDiff `json:"content"`
}

func (note *Note) getFormattedTitle() string {
	// TODO remove symbols and cap length
	line := strings.Split(note.Content, "\n")[0]
	trimmedLine := strings.TrimSpace(strings.TrimPrefix(line, "#"))
	return strings.ReplaceAll(strings.ToLower(trimmedLine), " ", "-")
}

func getNoteName(noteID string, note *Note) string {
	return note.getFormattedTitle() + "-" + noteID
}

func getFormattedTitle(summary *s.EntitySummary[Note]) string {
	line := strings.Split(summary.Data.Content, "\n")[0]
	trimmedLine := strings.TrimSpace(strings.TrimPrefix(line, "#"))
	return strings.ReplaceAll(strings.ToLower(trimmedLine), " ", "-")
}

func (client *client) getFileName(noteName string) string {
	return client.projectDir + "/" + noteName + ".gmi"
}

func (client *client) writeNote(noteID string, note *Note) error {
	name := getNoteName(noteID, note)
	filename := client.getFileName(name)

	if err := os.WriteFile(filename, []byte(note.Content), 0600); err != nil {
		return err
	}

	return nil
}
