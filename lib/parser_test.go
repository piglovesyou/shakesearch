package lib

import (
	"github.com/bradleyjkemp/cupaloy"
	"strings"
	"testing"
)

func TestParseTitles(t *testing.T) {
	content, _ := LoadContent("completeworks.txt")
	titles := ParseTitles(content, 0, &[]string{})
	cupaloy.SnapshotT(t, titles)
}

func TestBuildWorksIndex(t *testing.T) {
	content, _ := LoadContent("../completeworks.txt")
	works := BuildWorksIndex(content)
	for _, work := range *works {
		subContent := string((*content)[work.start : work.start+100])
		if strings.Contains(subContent, work.nameOnWork) == false {
			t.Fatal("boom")
		}
	}
	cupaloy.SnapshotT(t, works)
}
