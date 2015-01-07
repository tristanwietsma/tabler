package main

import (
	"log"
	"os"
	"regexp"
	"strings"
	"text/template"
)

var (
	filters = template.FuncMap{
		"plus1":  func(x int) int { return x + 1 },
		"lower":  func(s string) string { return strings.ToLower(s) },
		"caller": func(s string) string { return strings.ToLower(s)[0:1] },
	}
	tagPattern = regexp.MustCompile(`tabler:"([0-9a-zA-Z=&\(\)]*)"`)
)

func newTmpl(s string) *template.Template {
	return template.Must(template.New("T").Funcs(filters).Parse(s))
}

func main() {
	for _, path := range os.Args[1:] {
		infile := InputFile{}
		if err := infile.Init(path); err != nil {
			log.Printf("%v", err)
			continue
		}
		infile.Write()
	}
}
