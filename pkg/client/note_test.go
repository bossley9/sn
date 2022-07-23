package client

import (
	"testing"

	th "git.sr.ht/~bossley9/sn/pkg/testHelpers"
)

// GetContentTitleID

func TestGetContentTitleID_StandardHeading(t *testing.T) {
	test := `# Title Here
this is some content
another line here
`
	ref := "title-here"

	th.AssertEqual(t, GetContentTitleID(test), ref)
}

func TestGetContentTitleID_SymbolHeading(t *testing.T) {
	test := `# Please! "parse" this properly, or else
	this is some content
	another line here
	`
	ref := "please-parse-this-properly-or-el"

	th.AssertEqual(t, GetContentTitleID(test), ref)
}

func TestGetContentTitleID_LongTitle(t *testing.T) {
	test := `this is an extremely long title and I can't be bothered to cut it shorter if I'm being honest
	this is some content
	another line here
	`
	ref := "this-is-an-extremely-long-title"

	th.AssertEqual(t, GetContentTitleID(test), ref)
}

func TestGetContentTitleID_EndingHyphen(t *testing.T) {
	test := `1234567890123456789012345678901 hello
	this is some content
	another line here
	`
	ref := "1234567890123456789012345678901"

	th.AssertEqual(t, GetContentTitleID(test), ref)
}
