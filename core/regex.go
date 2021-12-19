package core

import (
	"regexp"
	"strings"
)

var ASCIIRegex = regexp.MustCompile("[[:^ascii:]]")
var WhitespaceRegex = regexp.MustCompile("\\s+")

func PrepareStringToBeUsedForHTMLID(text string) string {
	text = ASCIIRegex.ReplaceAllLiteralString(text, "")
	if len(text) > 30 {
		text = text[:30]
	}
	text = strings.Replace(strings.ToLower(text), " ", "_", -1)
	text = strings.Replace(strings.ToLower(text), ".", "_", -1)
	return text
}
