package lib

import (
	"bytes"
	"text/template"
)

// InsertRowStatement ...
func (t Table) InsertRowStatement() string {
	buf := bytes.Buffer{}
	tmpl := newTmpl(`INSERT INTO {{lower .Name}} ({{$n := len .Columns}}{{range $i, $c := .Columns}}{{lower $c.Name}}{{if lt (plus1 $i) $n}}, {{end}}{{end}}) VALUES ({{range $i, $c := .Columns}}?{{if lt (plus1 $i) $n}}, {{end}}{{end}});`)
	tmpl.Execute(&buf, t)
	return buf.String()
}

// InsertRow ...
func (t Table) InsertRow() string {
	buf := bytes.Buffer{}
	var tmpl *template.Template
	if !t.HasConn {
		tmpl = newTmpl(`func ({{caller .Name}} {{.Name}}) InsertRow() string {
    return ` + "`{{.InsertRowStatement}}`" + `
}`)
	} else {
		tmpl = newTmpl(`func ({{caller .Name}} {{.Name}}) InsertRow({{$n := len .Columns}}{{range $i, $c := .Columns}}{{lower $c.Name}} {{$c.FType}}{{if lt (plus1 $i) $n}}, {{end}}{{end}}) error {
    _, err := {{caller .Name}}.db.Exec(` + "`{{.InsertRowStatement}}`, " + `{{range $i, $c := .Columns}}{{lower $c.Name}}{{if lt (plus1 $i) $n}}, {{end}}{{end}})
    return err
}`)
	}
	tmpl.Execute(&buf, t)
	return buf.String()
}
