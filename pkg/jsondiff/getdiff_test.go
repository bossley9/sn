package jsondiff

import (
	"testing"

	th "git.sr.ht/~bossley9/gem/pkg/testhelpers"
)

func TestGetDiff_Routine(t *testing.T) {
	s1 := "Ted"
	s2 := "Red"
	ref := "-1\t+R\t=2"

	test := GetDiff(s1, s2)

	th.AssertEqual(t, test.Value, ref)
}

func TestGetDiff_Multiple(t *testing.T) {
	s1 := "the big cat walked to the store and ate"
	s2 := "the big dog walked to the mall and ate"
	// ref := "=8\t-3\t+dog\t=15\t-5\t+mall\t=8"
	// TEMP using swap and replace
	ref := "=8\t-23\t+dog walked to the mall\t=8"

	test := GetDiff(s1, s2)

	th.AssertEqual(t, test.Value, ref)
}

func TestGetDiff_Center(t *testing.T) {
	s1 := "the big cat walked to the store and ate"
	s2 := "walked"
	ref := "-12\t=6\t-21"

	test := GetDiff(s1, s2)

	th.AssertEqual(t, test.Value, ref)
}

func TestGetDiff_CenterAlt(t *testing.T) {
	s1 := "walked"
	s2 := "the big cat walked to the store and ate"
	ref := "+the big cat \t=6\t+ to the store and ate"

	test := GetDiff(s1, s2)

	th.AssertEqual(t, test.Value, ref)
}

func TestGetDiff_Same(t *testing.T) {
	s1 := "the big cat walked to the store and ate"
	s2 := "the big cat walked to the store and ate"
	ref := ""

	test := GetDiff(s1, s2)

	th.AssertEqual(t, test.Value, ref)
}
