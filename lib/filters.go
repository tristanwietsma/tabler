package lib

import (
	"strings"
	"text/template"
)

var (
	filters = template.FuncMap{
		"plus1":  func(x int) int { return x + 1 },
		"lower":  func(s string) string { return strings.ToLower(s) },
		"caller": func(s string) string { return strings.ToLower(s)[0:1] },
	}
)
