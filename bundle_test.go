package bundle_test

import "testing"

var data_dir = "test_data/"
var entries = []string{
	"car-sw.jpg", "don+peter.jpeg", "flowchart-woman.jpg",
}

func TestIndex(t *testing.T) {
	if sz := len(_bundleIdx); sz != len(entries) {
		t.Fatalf("Index size %d != %d", sz, len(entries))
	}
	if sz := len(_bundleIdx.Dir("")); sz != len(entries) {
		t.Fatalf("Dir returned %d != %d", sz, len(entries))
	}
	for _, e := range entries {
		if !_bundleIdx.Has(e) {
			t.Fatalf("Not in index: %s", e)
		}
	}
}
