package lib

import (
	"github.com/bradleyjkemp/cupaloy"
	"testing"
)

var fixture = `To be, or not to be, that is the question`

func TestCreateNGrams(t *testing.T) {
	tokens := Tokenize(fixture)
	unigrams := CreateNGrams(tokens, 1)
	bigrams := CreateNGrams(tokens, 2)
	trigrams := CreateNGrams(tokens, 3)
	cupaloy.SnapshotT(t, unigrams, bigrams, trigrams)
}
