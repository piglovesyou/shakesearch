package lib

import (
	"fmt"
	"regexp"
	"strings"
)

type Parser struct {
	CompleteWorks string
}

var contentsStartText = `\n      Contents\n`
var contentsStartRegexp = regexp.MustCompile(contentsStartText)
var contentsTitleRegexp = regexp.MustCompile(`^\n+ +\w[^\n]+\n`)
var omittableTitlePrefix = "THE TRAGEDY OF "
var dramatisPersonaeRegexp = regexp.MustCompile(`\nDramatis Person(ae|æ)\.?\n`)

func ParseTitles(content *string, pointer int, titles *[]string) *[]string {
	if pointer <= 0 {
		loc := contentsStartRegexp.FindStringIndex(*content)
		return ParseTitles(content, loc[1], titles)
	}
	loc := contentsTitleRegexp.FindStringIndex((*content)[pointer:])
	if loc == nil {
		return titles
	}
	start := pointer + loc[0]
	end := pointer + loc[1]
	*titles = append(*titles, Trim((*content)[start:end]))
	return ParseTitles(content, end, titles)
}

func findWorkLoc(content *string, pointer int, t string) []int {
	regexText := fmt.Sprintf("\n ?%v:?\n", t)
	r, _ := regexp.Compile(regexText)
	loc := r.FindStringIndex((*content)[pointer:])
	if loc != nil {
		return loc
	}

	// XXX: ugly patch 1.
	if strings.HasPrefix(t, omittableTitlePrefix) {
		return findWorkLoc(content, pointer, strings.TrimPrefix(t, omittableTitlePrefix))
	}

	// XXX: ugly patch 2.
	regexText = strings.ReplaceAll(regexText, ";", ":")
	r, _ = regexp.Compile(regexText)
	loc = r.FindStringIndex((*content)[pointer:])
	if loc != nil {
		return loc
	}

	// XXX: ugly patch 3.
	loc = dramatisPersonaeRegexp.FindStringIndex((*content)[pointer:])
	if loc != nil {
		return loc
	}

	return nil
}

func ParseWorks(content *string, pointer int, titles *[]string, works *[]Work) *[]Work {
	if len(*titles) <= 0 {
		return works
	}
	t := (*titles)[0]
	loc := findWorkLoc(content, pointer, t)
	if loc == nil {
		return works
	}
	start := pointer + loc[0]
	end := pointer + loc[1]
	*works = append(*works, Work{name: t, start: start, nameOnWork: Trim((*content)[start:end])})
	tRest := (*titles)[1:]
	return ParseWorks(content, end, &tRest, works)
}

func BuildWorksIndex(content *string) *[]Work {
	titles := ParseTitles(content, 0, &[]string{})
	works := ParseWorks(content, 0, titles, &[]Work{})
	return works
}
