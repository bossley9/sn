package jsondiff

import (
	"testing"

	th "git.sr.ht/~bossley9/gem/pkg/testhelpers"
)

func TestApply_Routine(t *testing.T) {
	test := "Ted"
	ref := "Red"

	diff := StringJSONDiff{
		Operation: "d",
		Value:     "-1\t+R\t=2",
	}

	th.AssertEqual(t, diff.Apply(test), ref)
}

func TestApply_Multiple(t *testing.T) {
	test := "the big cat walked to the store and ate"
	ref := "the big dog walked to the mall and ate"

	diff := StringJSONDiff{
		Operation: "d",
		Value:     "=8\t-3\t+dog\t=15\t-5\t+mall\t=8",
	}

	th.AssertEqual(t, diff.Apply(test), ref)
}
