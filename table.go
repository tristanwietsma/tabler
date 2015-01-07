package main

import (
	"bytes"
)

// Table is an SQL table.
type Table struct {
	Name        string
	HasConn     bool
	Conn        string
	Columns     []Column
	PrimaryKeys []Column
}

func (t Table) CreateTableStatement() string {
	buf := bytes.Buffer{}
	tmpl := newTmpl(`CREATE TABLE {{lower .Name}} ({{$n := len .Columns}}{{range $i, $c := .Columns}}{{$c.String}}{{if lt (plus1 $i) $n}}, {{end}}{{end}}){{$p := len .PrimaryKeys}}{{if $p}} PRIMARY KEY ({{range $j, $k := .PrimaryKeys}}{{lower $k.Name}}{{if lt (plus1 $j) $p}}, {{end}}{{end}}){{end}};`)
	tmpl.Execute(&buf, t)
	return buf.String()
}

func (t Table) CreateTable() string {
	buf := bytes.Buffer{}
	tmpl := newTmpl(`func ({{caller .Name}} {{.Name}}) CreateTable() string {
    return ` + "`" + `{{.CreateTableStatement}}` + "`" + `
}`)
	tmpl.Execute(&buf, t)
	return buf.String()
}

func (t Table) DropTable() string {
	buf := bytes.Buffer{}
	tmpl := newTmpl(`func ({{caller .Name}} {{.Name}}) DropTable() string {
    return ` + "`" + `DROP TABLE {{lower .Name}};` + "`" + `
}`)
	tmpl.Execute(&buf, t)
	return buf.String()
}

func (t Table) InsertRow() string {
	buf := bytes.Buffer{}
	tmpl := newTmpl(`func ({{caller .Name}} {{.Name}}) InsertRow() string {
    return ` + "`" + `INSERT INTO {{lower .Name}} ({{$n := len .Columns}}{{range $i, $c := .Columns}}{{lower $c.Name}}{{if lt (plus1 $i) $n}}, {{end}}{{end}}) VALUES ({{range $i, $c := .Columns}}?{{if lt (plus1 $i) $n}}, {{end}}{{end}});` + "`" + `
}`)
	tmpl.Execute(&buf, t)
	return buf.String()
}

func (t Table) SelectRow() string {
	buf := bytes.Buffer{}
	tmpl := newTmpl(`func ({{caller .Name}} {{.Name}}) SelectRow() string {
    return ` + "`" + `SELECT {{$n := len .Columns}}{{range $i, $c := .Columns}}{{lower $c.Name}}{{if lt (plus1 $i) $n}}, {{end}}{{end}} FROM {{lower .Name}} WHERE {{$p := len .PrimaryKeys}}{{range $j, $k := .PrimaryKeys}}{{lower $k.Name}}=?{{if lt (plus1 $j) $p}} AND {{end}}{{end}};` + "`" + `
}`)
	tmpl.Execute(&buf, t)
	return buf.String()
}
