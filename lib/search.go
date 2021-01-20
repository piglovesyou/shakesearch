package lib

import (
	"fmt"
	"io/ioutil"
	"math"
	"sort"
)

type Searcher struct {
	CompleteWorks *string
	WorksIndex    *[]Work
	NGramRules    *[]*NGramRule
	//SuffixArray *suffixarray.Index
}

type SearchResultHighlight struct {
	SubContent string
	Start      int
	End        int
}

type SearchResult struct {
	Name       string
	Highlights []SearchResultHighlight
	Score      int
}

type SearchResponse struct {
	query    string
	timeTook int
	results  SearchResult
}

func (s *Searcher) Load(filename string) error {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Load: %w", err)
	}
	completeWorks := string(CrossPlatformNewlineRegexp.ReplaceAll(dat, []byte("\n")))
	s.CompleteWorks = &completeWorks

	worksIndex := BuildWorksIndex(&completeWorks)
	s.WorksIndex = worksIndex

	nGramRules := []*NGramRule{}
	nGramRules = append(nGramRules, &NGramRule{"trigram", &NGramMap{}, 3, 4, nil})
	nGramRules = append(nGramRules, &NGramRule{"bigram", &NGramMap{}, 2, 2, nil})
	nGramRules = append(nGramRules, &NGramRule{"unigram", &NGramMap{}, 1, 1, nil})
	for wIndex, work := range *worksIndex {
		wStart := work.start
		var wEnd int
		if wIndex+1 < len(*worksIndex) {
			wEnd = (*worksIndex)[wIndex+1].start
		} else {
			wEnd = len(completeWorks)
		}
		wContent := completeWorks[wStart:wEnd]
		tokens := *TokenizeWithIndex(wStart, &wContent)

		for ti, _ := range tokens {
			for _, rule := range nGramRules {
				rule.AppendToken(wIndex, ti, &tokens)
			}
		}
	}
	for _, rule := range nGramRules {
		rule.Finalize()
	}
	s.NGramRules = &nGramRules

	//s.SuffixArray = suffixarray.New(dat)
	return nil
}

type SearchResultMapValue struct {
	score int
	poss  *[]NGramPos
}
type SearchResultMap = map[int]*SearchResultMapValue

var SearchResultContentLen = 300

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (s *Searcher) FinalizeSearchResults(m *SearchResultMap) *[]SearchResult {
	results := []SearchResult{}
	for workIndex, searchResult := range *m {
		poss := searchResult.poss
		work := (*s.WorksIndex)[workIndex]
		offsetLen := int(math.Floor(float64(SearchResultContentLen) / float64(len(*poss)) / 2))
		highlights := []SearchResultHighlight{}
		for _, pos := range *poss {
			offsetStart := Max(pos.start-offsetLen, 0)
			offsetEnd := Min(pos.end+offsetLen, len(*s.CompleteWorks))
			subContent := string((*s.CompleteWorks)[offsetStart:offsetEnd])
			highlights = append(highlights, SearchResultHighlight{
				SubContent: subContent,
				Start:      offsetLen,
				End:        offsetLen + (pos.end - pos.start),
			})
		}
		results = append(results, SearchResult{
			Name:       work.name,
			Highlights: highlights,
			Score:      searchResult.score,
		})
	}
	sort.Slice(results, func(a, b int) bool {
		// Order by desc
		return results[a].Score > results[b].Score
	})
	return &results
}

var maxCandidateLen = 5

func hasIntersection(p1, p2 NGramPos) bool {
	return (p1.start <= p2.start && p2.start <= p1.end) ||
		(p2.start <= p1.start && p1.start <= p2.end)
}

func AppendOrMerge(poss *[]NGramPos, pos NGramPos) *[]NGramPos {
	for i, p := range *poss {
		if hasIntersection(p, pos) {
			(*poss)[i] = NGramPos{
				workIndex: p.workIndex,
				start:     Min(p.start, pos.start),
				end:       Max(p.end, pos.end),
			}
			return poss
		}
	}
	*poss = append(*poss, pos)
	return poss
}

func AppendResultPoss(resultsMap *SearchResultMap, queryTokens *[]string, ruleP *NGramRule) *SearchResultMap {
	rule := *ruleP
	for _, qt := range CreateNGrams(*queryTokens, rule.n) {
		for _, pos := range *rule.SearchToken(qt) {
			if searchResult, ok := (*resultsMap)[pos.workIndex]; ok {
				poss := searchResult.poss
				searchResult.score = searchResult.score + rule.score
				poss = AppendOrMerge(poss, pos)
				if maxCandidateLen <= len(*poss) {
					break
				}
			} else {
				poss := []NGramPos{pos}
				(*resultsMap)[pos.workIndex] = &SearchResultMapValue{rule.score, &poss}
			}
		}
	}
	return resultsMap
}

func (s *Searcher) Search(query string) *[]SearchResult {
	queryTokens := Tokenize(query)
	resultsMap := SearchResultMap{}

	for _, rule := range *s.NGramRules {
		AppendResultPoss(&resultsMap, &queryTokens, rule)
	}

	//idxs := s.SuffixArray.Lookup([]byte(query), -1)
	//for _, idx := range idxs {
	//	results = append(results, s.CompleteWorks[idx-250:idx+250])
	//}

	return s.FinalizeSearchResults(&resultsMap)
}
