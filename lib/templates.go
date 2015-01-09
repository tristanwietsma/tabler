package lib

import (
	"bytes"
	"text/template"
)

func newTmpl(s string) *template.Template {
	return template.Must(template.New("T").Funcs(filters).Parse(s))
}
