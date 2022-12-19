package client

import (
	"errors"
	"os"

	"git.sr.ht/~bossley9/gem/pkg/url"
	f "git.sr.ht/~bossley9/sn/pkg/fileio"
	j "git.sr.ht/~bossley9/sn/pkg/jsondiff"
)

type NoteResponse struct {
	Tags           []string `json:"tags"`
	Deleted        bool     `json:"deleted"`
	ShareURL       string   `json:"shareURL"`
	PublishURL     string   `json:"publishURL"`
	Content        string   `json:"content"`
	SystemTags     []string `json:"systemTags"`
	LastEditedDate float32  `json:"modificationDate"`
	CreationDate   float32  `json:"creationDate"`
}

type NoteID string

type Note struct {
	// do not store ID, it will be used as the map key
	Version int    `json:"v"`
	Name    string `json:"n"`
}

type NoteSummary struct {
	ID      NoteID
	Version int
	Content string
}

type NoteDiff struct {
	Content      j.StringJSONDiff  `json:"content"`
	Deleted      j.BoolJSONDiff    `json:"deleted"`
	CreationDate j.Float32JSONDiff `json:"creationDate"`
}

// given a note id and content string, returns a unique note name identifier
func GetNoteName(noteID NoteID, content string) string {
	return url.GenerateID(content) + "-" + string(noteID)
}

// given a note name, returns an absolute path filename
func (client *Client) getFileName(noteName string) string {
	return client.projectDir + "/" + noteName + ".md"
}

// given a note name, returns an absolute path version filename
func (client *Client) getVersionFileName(noteName string) string {
	return client.versionDir + "/" + noteName + ".md"
}

// given a note summary, writes the note to file and updates the cache and version if necessary
func (client *Client) writeNote(summary *NoteSummary) error {
	noteName := ""
	note, ok := client.storage.Notes[summary.ID]
	if ok {
		noteName = note.Name
	} else {
		noteName = GetNoteName(summary.ID, summary.Content)
	}

	// write note to file
	filename := client.getFileName(noteName)
	if err := os.WriteFile(filename, []byte(summary.Content), f.RW); err != nil {
		return err
	}
	vFilename := client.getVersionFileName(noteName)
	if err := os.WriteFile(vFilename, []byte(summary.Content), f.RW); err != nil {
		return err
	}

	// update cache
	client.storage.Notes[summary.ID] = Note{
		Version: summary.Version,
		Name:    noteName,
	}

	return nil
}

// given a note id, returns written content associated with that note
func (client *Client) readNote(noteID NoteID) (string, error) {
	note, ok := client.storage.Notes[noteID]
	if !ok {
		return "", errors.New("note with id " + string(noteID) + " does not exist.")
	}

	filename := client.getFileName(note.Name)
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func (client *Client) readVersionNote(noteID NoteID) (string, error) {
	note, ok := client.storage.Notes[noteID]
	if !ok {
		return "", errors.New("note with id " + string(noteID) + " does not exist.")
	}

	filename := client.getVersionFileName(note.Name)
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}

	return string(content), nil
}
