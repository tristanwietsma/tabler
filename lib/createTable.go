package lib

import "bytes"

// CreateTableStatement ...
func (t Table) CreateTableStatement() string {
	buf := bytes.Buffer{}
	templify(`CREATE TABLE {{lower .Name}} ({{$n := len .Columns}}{{range $i, $c := .Columns}}{{$c.String}}{{if lt (plus1 $i) $n}}, {{end}}{{end}}){{$p := len .PrimaryKeys}}{{if $p}} PRIMARY KEY ({{range $j, $k := .PrimaryKeys}}{{lower $k.Name}}{{if lt (plus1 $j) $p}}, {{end}}{{end}}){{end}};`).Execute(&buf, t)
	return buf.String()
}

// CreateTable ...
func (t Table) CreateTable() string {
	buf := bytes.Buffer{}
	if !t.HasConn {
		templify(`func ({{caller .Name}} {{.Name}}) CreateTable() string {
    return "{{.CreateTableStatement}}"
}`).Execute(&buf, t)
	} else {
		templify(`func ({{caller .Name}} {{.Name}}) CreateTable() error {
    _, err := {{caller .Name}}.db.Exec("{{.CreateTableStatement}}")
    return err
}`).Execute(&buf, t)
	}
	return buf.String()
}
