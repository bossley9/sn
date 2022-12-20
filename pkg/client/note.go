package client

import (
	"errors"
	"os"

	"git.sr.ht/~bossley9/gem/pkg/url"
	f "git.sr.ht/~bossley9/sn/pkg/fileio"
	j "git.sr.ht/~bossley9/sn/pkg/jsondiff"
	s "git.sr.ht/~bossley9/sn/pkg/simperium"
)

type NoteResponse struct {
	Tags             []string `json:"tags"`
	Deleted          bool     `json:"deleted"`
	ShareURL         string   `json:"shareURL"`
	PublishURL       string   `json:"publishURL"`
	Content          string   `json:"content"`
	SystemTags       []string `json:"systemTags"`
	ModificationDate float32  `json:"modificationDate"`
	CreationDate     float32  `json:"creationDate"`
}

type NoteID string

type Note struct {
	// do not store ID, it will be used as the map key
	Version int    `json:"v"`
	Name    string `json:"n"`
}

type NoteDiff struct {
	Tags             j.JSONDiff[[]string] `json:"tags,omitempty"`
	Deleted          j.BoolJSONDiff       `json:"deleted,omitempty"`
	ShareURL         j.StringJSONDiff     `json:"shareURL,omitempty"`
	PublishURL       j.StringJSONDiff     `json:"publishURL,omitempty"`
	Content          j.StringJSONDiff     `json:"content,omitempty"`
	SystemTags       j.JSONDiff[[]string] `json:"systemTags,omitempty"`
	ModificationDate j.Int64JSONDiff      `json:"modificationDate,omitempty"`
	CreationDate     j.Float32JSONDiff    `json:"creationDate,omitempty"`
}

type DownloadNoteDiff s.Change[NoteDiff]
type UploadNoteDiff s.UploadChange[NoteDiff]

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

func (client *Client) writeNote(noteID NoteID, note *Note, content string) error {
	// write note to file
	filename := client.getFileName(note.Name)
	if err := os.WriteFile(filename, []byte(content), f.RW); err != nil {
		return err
	}
	vFilename := client.getVersionFileName(note.Name)
	if err := os.WriteFile(vFilename, []byte(content), f.RW); err != nil {
		return err
	}

	client.storage.Notes[noteID] = Note{
		Version: note.Version,
		Name:    note.Name,
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
