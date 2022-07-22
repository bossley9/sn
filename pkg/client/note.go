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

func getFormattedTitle(summary *s.EntitySummary[Note]) string {
	line := strings.Split(summary.Data.Content, "\n")[0]
	trimmedLine := strings.TrimSpace(strings.TrimPrefix(line, "#"))
	return strings.ReplaceAll(strings.ToLower(trimmedLine), " ", "-")
}

func (client *client) getFileName(summary *s.EntitySummary[Note]) string {
	root := client.projectDir
	title := getFormattedTitle(summary)
	return root + "/" + title + "-" + summary.ID + ".gmi"
}

func (client *client) getFileNameFromID(entityID string) (string, error) {
	files, err := os.ReadDir(client.projectDir)
	if err != nil {
		return "", err
	}

	for _, file := range files {
		filename := file.Name()
		if strings.Contains(filename, entityID) {
			return client.projectDir + "/" + filename, nil
		}
	}

	return "", nil
}

func (client *client) writeNoteSummary(summary *s.EntitySummary[Note]) error {
	filename := client.getFileName(summary)
	if err := os.WriteFile(filename, []byte(summary.Data.Content), 0600); err != nil {
		return err
	}

	return nil
}
