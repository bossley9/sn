package client

import (
	"os"
	"strings"
)

type Note struct {
	// returned from index command
	ID      string `json:"id"`
	Version int    `json:"version"`
	// returned from entity command
	Tags           []string `json:"tags"`
	Deleted        bool     `json:"deleted"`
	ShareURL       string   `json:"shareURL"`
	PublishURL     string   `json:"publishURL"`
	Content        string   `json:"content"`
	SystemTags     []string `json:"systemTags"`
	LastEditedDate float32  `json:"modificationDate"`
	CreationDate   float32  `json:"creationDate"`
}

func (note *Note) getFormattedTitle() string {
	line := strings.Split(note.Content, "\n")[0]
	trimmedLine := strings.TrimSpace(strings.TrimPrefix(line, "#"))

	return strings.ReplaceAll(strings.ToLower(trimmedLine), " ", "-")
}

func (client *client) writeNote(note *Note) error {
	root := client.projectDir
	title := note.getFormattedTitle()
	filename := root + "/" + title + "-" + note.ID + ".gmi"

	if err := os.WriteFile(filename, []byte(note.Content), 0600); err != nil {
		return err
	}

	return nil
}
