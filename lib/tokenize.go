package lib

import (
	"regexp"
	"strings"
)

var tokenSeparatorPatternSource = `[^A-Za-zА-Яа-я0-9_]+`
var tokenSeparatorRegExp = regexp.MustCompile(tokenSeparatorPatternSource)

func Tokenize(content string) []string {
	delim := " "
	tokenized := tokenSeparatorRegExp.ReplaceAll([]byte(content), []byte(delim))
	return strings.Split(SanitizeToken(string(tokenized)), delim)
}

type Token struct {
	value string
	start int
	end   int
}

func TokenizeWithIndex(workOffset int, contentP *string) *[]Token {
	content := *contentP
	tokens := []Token{}
	positions := tokenSeparatorRegExp.FindAllStringIndex(content, -1)
	offset := 0
	for _, pos := range positions {
		tokens = append(tokens, Token{
			SanitizeToken(content[offset:pos[0]]),
			workOffset + offset,
			workOffset + pos[0],
		})
		offset = pos[1]
	}
	tokens = append(tokens, Token{
		SanitizeToken(content[offset:]),
		workOffset + offset,
		workOffset + len(content),
	})
	return &tokens
}

func CreateNGrams(tokens []string, n int) []string {
	rv := []string{}
	for ti, _ := range tokens {
		if ti+n <= len(tokens) {
			ngramToken := []string{}
			for i := ti; i < ti+n; i++ {
				ngramToken = append(ngramToken, tokens[i])
			}
			rv = append(rv, strings.Join(ngramToken, " "))
		}
	}
	return rv
}
