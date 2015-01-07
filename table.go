package main

import (
	"bytes"
)

// Table is an SQL table.
type Table struct {
	Name        string
	Columns     []Column
	PrimaryKeys []Column
}

// CreateTable returns a create table statement for the table.
func (t Table) CreateTable() string {
	buf := bytes.Buffer{}
	tmpl := newTmpl(`func ({{caller .Name}} {{.Name}}) CreateTable() string {
    return ` + "`" + `CREATE TABLE {{lower .Name}} ({{$n := len .Columns}}{{range $i, $c := .Columns}}{{$c.String}}{{if lt (plus1 $i) $n}}, {{end}}{{end}}){{$p := len .PrimaryKeys}}{{if $p}} PRIMARY KEY ({{range $j, $k := .PrimaryKeys}}{{lower $k.Name}}{{if lt (plus1 $j) $p}}, {{end}}{{end}}){{end}};` + "`" + `
}`)
	tmpl.Execute(&buf, t)
	return buf.String()
}

// DropTable returns a drop table statement for the table.
func (t Table) DropTable() string {
	buf := bytes.Buffer{}
	tmpl := newTmpl(`func ({{caller .Name}} {{.Name}}) DropTable() string {
    return ` + "`" + `DROP TABLE {{lower .Name}};` + "`" + `
}`)
	tmpl.Execute(&buf, t)
	return buf.String()
}

// InsertRow returns a parameterized insertion statement for the table.
func (t Table) InsertRow() string {
	buf := bytes.Buffer{}
	tmpl := newTmpl(`func ({{caller .Name}} {{.Name}}) InsertRow() string {
    return ` + "`" + `INSERT INTO {{lower .Name}} ({{$n := len .Columns}}{{range $i, $c := .Columns}}{{lower $c.Name}}{{if lt (plus1 $i) $n}}, {{end}}{{end}}) VALUES ({{range $i, $c := .Columns}}?{{if lt (plus1 $i) $n}}, {{end}}{{end}});` + "`" + `
}`)
	tmpl.Execute(&buf, t)
	return buf.String()
}

// SelectRow returns a parameterized query for a single row.
func (t Table) SelectRow() string {
	buf := bytes.Buffer{}
	tmpl := newTmpl(`func ({{caller .Name}} {{.Name}}) SelectRow() string {
    return ` + "`" + `SELECT {{$n := len .Columns}}{{range $i, $c := .Columns}}{{lower $c.Name}}{{if lt (plus1 $i) $n}}, {{end}}{{end}} FROM {{lower .Name}} WHERE {{$p := len .PrimaryKeys}}{{range $j, $k := .PrimaryKeys}}{{lower $k.Name}}=?{{if lt (plus1 $j) $p}} AND {{end}}{{end}};` + "`" + `
}`)
	tmpl.Execute(&buf, t)
	return buf.String()
}
