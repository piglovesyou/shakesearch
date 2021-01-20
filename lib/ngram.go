package lib

import (
	"sort"
	"strings"
)

type NGramPos struct {
	workIndex int
	start     int
	end       int
}

type NGramMap = map[string]*[]NGramPos

type NGramRule struct {
	name         string
	nGramMap     *NGramMap
	n            int
	score        int
	sortedTokens *[]string
}

func (g *NGramRule) AppendToken(workIndex int, tokenIndex int, tokensP *[]Token) {
	tokens := *tokensP
	ngramMap := *g.nGramMap
	if len(tokens[tokenIndex].value) <= 0 {
		return
	}
	if tokenIndex+g.n <= len(tokens) {
		ngramTokens := tokens[tokenIndex : tokenIndex+g.n]
		tokenArray := []string{}
		for _, ngramToken := range ngramTokens {
			tokenArray = append(tokenArray, ngramToken.value)
		}
		t := Trim(strings.Join(tokenArray, " "))
		if len(t) > 0 {
			start := ngramTokens[0].start
			end := ngramTokens[g.n-1].end
			poss := ngramMap[t]
			if poss == nil {
				poss = &[]NGramPos{}
				ngramMap[t] = poss
			}
			*poss = append(*poss, NGramPos{workIndex, start, end})
		}
	}
}

func (g *NGramRule) Finalize() {
	ngramMap := *g.nGramMap
	ngramTokens := []string{}
	for t, _ := range ngramMap {
		ngramTokens = append(ngramTokens, t)
	}
	sort.Slice(ngramTokens, func(a, b int) bool {
		return ngramTokens[a] < ngramTokens[b]
	})
	g.sortedTokens = &ngramTokens
}

func (g *NGramRule) SearchToken(token string) *[]NGramPos {
	nGramMap := *g.nGramMap
	sortedTokens := *g.sortedTokens
	i := sort.SearchStrings(sortedTokens, token)
	return nGramMap[sortedTokens[i]]
}
