package main

import (
	"bytes"
	"text/template"
)

// Table is an SQL table.
type Table struct {
	Name        string
	HasConn     bool
	Conn        string
	Columns     []Column
	PrimaryKeys []Column
}

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

// DropTableStatement ...
func (t Table) DropTableStatement() string {
	buf := bytes.Buffer{}
	tmpl := newTmpl(`DROP TABLE {{lower .Name}};`)
	tmpl.Execute(&buf, t)
	return buf.String()
}

// DropTable ...
func (t Table) DropTable() string {
	buf := bytes.Buffer{}
	var tmpl *template.Template
	if !t.HasConn {
		tmpl = newTmpl(`func ({{caller .Name}} {{.Name}}) DropTable() string {
    return ` + "`{{.DropTableStatement}}`" + `
}`)
	} else {
		tmpl = newTmpl(`func ({{caller .Name}} {{.Name}}) DropTable() error {
    _, err := {{caller .Name}}.db.Exec(` + "`{{.DropTableStatement}}`)" + `
    return err
}`)
	}
	tmpl.Execute(&buf, t)
	return buf.String()
}

// InsertRow ...
func (t Table) InsertRow() string {
	buf := bytes.Buffer{}
	tmpl := newTmpl(`func ({{caller .Name}} {{.Name}}) InsertRow() string {
    return ` + "`" + `INSERT INTO {{lower .Name}} ({{$n := len .Columns}}{{range $i, $c := .Columns}}{{lower $c.Name}}{{if lt (plus1 $i) $n}}, {{end}}{{end}}) VALUES ({{range $i, $c := .Columns}}?{{if lt (plus1 $i) $n}}, {{end}}{{end}});` + "`" + `
}`)
	tmpl.Execute(&buf, t)
	return buf.String()
}

// SelectRow ...
func (t Table) SelectRow() string {
	buf := bytes.Buffer{}
	tmpl := newTmpl(`func ({{caller .Name}} {{.Name}}) SelectRow() string {
    return ` + "`" + `SELECT {{$n := len .Columns}}{{range $i, $c := .Columns}}{{lower $c.Name}}{{if lt (plus1 $i) $n}}, {{end}}{{end}} FROM {{lower .Name}} WHERE {{$p := len .PrimaryKeys}}{{range $j, $k := .PrimaryKeys}}{{lower $k.Name}}=?{{if lt (plus1 $j) $p}} AND {{end}}{{end}};` + "`" + `
}`)
	tmpl.Execute(&buf, t)
	return buf.String()
}
