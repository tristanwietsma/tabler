package lib

import "text/template"

func templify(s string) *template.Template {
	return template.Must(template.New("T").Funcs(filters).Parse(s))
}
