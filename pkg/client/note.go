package client

import (
	"os"
	"strings"

	j "git.sr.ht/~bossley9/sn/pkg/jsondiff"
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

type NoteDiff struct {
	Content j.StringJSONDiff `json:"content"`
}

func (note *Note) getFormattedTitle() string {
	line := strings.Split(note.Content, "\n")[0]
	trimmedLine := strings.TrimSpace(strings.TrimPrefix(line, "#"))

	return strings.ReplaceAll(strings.ToLower(trimmedLine), " ", "-")
}

func (client *client) getFileName(note *Note) string {
	root := client.projectDir
	title := note.getFormattedTitle()
	return root + "/" + title + "-" + note.ID + ".gmi"
}

func (client *client) getFileNameFromID(noteID string) (string, error) {
	files, err := os.ReadDir(client.projectDir)
	if err != nil {
		return "", err
	}

	for _, file := range files {
		filename := file.Name()
		if strings.Contains(filename, noteID) {
			return client.projectDir + "/" + filename, nil
		}
	}

	return "", nil
}

func (client *client) writeNote(note *Note) error {
	filename := client.getFileName(note)
	if err := os.WriteFile(filename, []byte(note.Content), 0600); err != nil {
		return err
	}

	return nil
}
