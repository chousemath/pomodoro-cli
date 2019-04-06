package stredit

import "testing"

func TestPluralize(t *testing.T) {
	s := Pluralize(2)
	if s != "s" {
		t.Error("Expected pluralization to happen, but it didn't")
	}
}
