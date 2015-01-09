package lib

import (
	"bytes"
	"text/template"
)

// CreateTableStatement ...
func (t Table) CreateTableStatement() string {
	buf := bytes.Buffer{}
	tmpl := newTmpl(`CREATE TABLE {{lower .Name}} ({{$n := len .Columns}}{{range $i, $c := .Columns}}{{$c.String}}{{if lt (plus1 $i) $n}}, {{end}}{{end}}){{$p := len .PrimaryKeys}}{{if $p}} PRIMARY KEY ({{range $j, $k := .PrimaryKeys}}{{lower $k.Name}}{{if lt (plus1 $j) $p}}, {{end}}{{end}}){{end}};`)
	tmpl.Execute(&buf, t)
	return buf.String()
}

// CreateTable ...
func (t Table) CreateTable() string {
	buf := bytes.Buffer{}
	var tmpl *template.Template
	if !t.HasConn {
		tmpl = newTmpl(`func ({{caller .Name}} {{.Name}}) CreateTable() string {
    return ` + "`" + `{{.CreateTableStatement}}` + "`" + `
}`)
	} else {
		tmpl = newTmpl(`func ({{caller .Name}} {{.Name}}) CreateTable() error {
    _, err := {{caller .Name}}.db.Exec(` + "`{{.CreateTableStatement}}`)" + `
    return err
}`)
	}
	tmpl.Execute(&buf, t)
	return buf.String()
}
