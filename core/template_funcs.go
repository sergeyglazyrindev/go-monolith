package core

import (
	"fmt"
	"html/template"
	"strings"
	"unicode"
	"unicode/utf8"
)

func add(n1 int, n2 int) int {
	return n1 + n2
}

func mul(n1 int, n2 int) int {
	return n1 * n2
}

func safe(s string) template.HTML {
	return template.HTML(s)
}

func GenerateAttrs(attrs map[string]string) template.HTML {
	attrsContent := make([]string, 0)
	for k, v := range attrs {
		attrsContent = append(attrsContent, fmt.Sprintf(" %s=\"%s\" ", template.HTMLAttr(k), template.HTML(v)))
	}
	return template.HTML(strings.Join(attrsContent, " "))
}

func attr(s string) template.HTMLAttr {
	return template.HTMLAttr(s)
}

func GetDisplayName(src string) string {
	// don't split invalid utf8
	if !utf8.ValidString(src) {
		return src
	}
	entries := []string{}
	var runes [][]rune
	lastClass := 0
	class := 0
	// split into fields based on class of unicode character
	for _, r := range src {
		switch true {
		case unicode.IsLower(r):
			class = 1
		case unicode.IsUpper(r):
			class = 2
		case unicode.IsDigit(r):
			class = 3
		default:
			class = 4
		}
		if class == lastClass {
			runes[len(runes)-1] = append(runes[len(runes)-1], r)
		} else {
			runes = append(runes, []rune{r})
		}
		lastClass = class
	}
	// handle upper case -> lower case sequences, e.g.
	// "PDFL", "oader" -> "PDF", "Loader"
	for i := 0; i < len(runes)-1; i++ {
		if unicode.IsUpper(runes[i][0]) && unicode.IsLower(runes[i+1][0]) {
			runes[i+1] = append([]rune{runes[i][len(runes[i])-1]}, runes[i+1]...)
			runes[i] = runes[i][:len(runes[i])-1]
		}
	}
	// construct []string from results
	for _, s := range runes {
		if len(s) > 0 {
			entries = append(entries, string(s))
		}
	}
	return strings.Join(entries, " ")
}

var FuncMap = template.FuncMap{
	"Tf":             Tf,
	"add":            add,
	"mul":            mul,
	"safe":           safe,
	"attr":           attr,
	"GenerateAttrs":  GenerateAttrs,
	"GetDisplayName": GetDisplayName,
	"SplitCamelCase": HumanizeCamelCase,
	//"Translate": func (v interface{}) string {
	//	return v.(string)
	//},
	//"CSRF": func() string {
	//	return "dfsafsa"
	//	// return authapi.GetSession(r)
	//},
}
