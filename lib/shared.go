package lib

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
)

var CrossPlatformNewlineRegexp = regexp.MustCompile(`\r?\n`)

type Work struct {
	name       string
	nameOnWork string
	start      int
}

func LoadContent(filename string) (*string, error) {
	bin, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("Load: %w", err)
	}
	content := string(CrossPlatformNewlineRegexp.ReplaceAll(bin, []byte("\n")))
	return &content, nil
}

func Trim(s string) string {
	return strings.Trim(s, "\n ")
}

func SanitizeToken(t string) string {
	return strings.ToLower(strings.ReplaceAll(t, "_", ""))
}
