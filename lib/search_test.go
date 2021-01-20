package lib

import (
	"github.com/bradleyjkemp/cupaloy"
	"testing"
)

func TestSearch(t *testing.T) {
	s := Searcher{}

	err := s.Load("completeworks.txt")
	if err != nil {
		t.Fatal(err)
	}

	results := s.Search("hamlet to be ore")
	cupaloy.SnapshotT(t, results)
}
