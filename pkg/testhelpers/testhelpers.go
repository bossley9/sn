package testhelpers

import "testing"

func AssertEqual[T comparable](t *testing.T, test, ref T) {
	t.Helper()

	if test != ref {
		t.Errorf("Expected:\n%v\nbut received:\n%v", ref, test)
	}
}
