package client

import (
	"log"
	"os"
	"regexp"
	"strings"

	j "git.sr.ht/~bossley9/sn/pkg/jsondiff"
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

type NoteSummary struct {
	ID      string
	Version int
	Content string
}

type NoteDiff struct {
	Content j.StringJSONDiff `json:"content"`
}

// given a content string of text, returns a formatted title in the form of an ID
// 1. get first line
// 2. remove heading indicator (#) and trim whitespace
// 3. convert to lowercase and replace whitespace with hyphens
// 4. remove symbols
// 5. cap length to maxLen chars and trim hyphen suffix
func GetContentTitleID(content string) string {
	maxLen := 32
	r, err := regexp.CompilePOSIX("[^a-zA-Z0-9-]+")
	if err != nil {
		log.Fatal("unable to parse title id regular expression. Exiting.")
	}

	firstLine := strings.Split(content, "\n")[0]
	trimmedLine := strings.TrimSpace(strings.TrimPrefix(firstLine, "#"))
	lowerLine := strings.ReplaceAll(strings.ToLower(trimmedLine), " ", "-")
	sanitizedLine := r.ReplaceAllString(lowerLine, "")

	cappedLine := sanitizedLine
	if len(cappedLine) > maxLen {
		cappedLine = cappedLine[:maxLen]
	}

	return strings.TrimSuffix(cappedLine, "-")
}

// given a note id and content string, returns a unique note name identifier
func GetNoteName(noteID string, content string) string {
	return GetContentTitleID(content) + "-" + noteID
}

func (client *client) getFileName(noteName string) string {
	return client.projectDir + "/" + noteName + ".gmi"
}

// given a note summary, writes the note to file and updates the cache if necessary
func (client *client) writeNote(summary *NoteSummary) error {
	// check for note name from cache
	if client.cache.Notes == nil {
		client.cache.Notes = make(map[string]NoteCache)
	}
	noteName := ""
	noteCache, ok := client.cache.Notes[summary.ID]
	if ok {
		noteName = noteCache.Name
	} else {
		noteName = GetNoteName(summary.ID, summary.Content)
	}

	// write note to file
	filename := client.getFileName(noteName)
	if err := os.WriteFile(filename, []byte(summary.Content), 0600); err != nil {
		return err
	}

	// update cache
	client.cache.Notes[summary.ID] = NoteCache{
		Version: summary.Version,
		Name:    noteName,
	}

	return client.writeCache()
}
