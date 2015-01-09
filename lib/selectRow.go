package lib

import (
	"bytes"
)

// SelectRow ...
func (t Table) SelectRow() string {
	buf := bytes.Buffer{}
	tmpl := newTmpl(`func ({{caller .Name}} {{.Name}}) SelectRow() string {
    return ` + "`" + `SELECT {{$n := len .Columns}}{{range $i, $c := .Columns}}{{lower $c.Name}}{{if lt (plus1 $i) $n}}, {{end}}{{end}} FROM {{lower .Name}} WHERE {{$p := len .PrimaryKeys}}{{range $j, $k := .PrimaryKeys}}{{lower $k.Name}}=?{{if lt (plus1 $j) $p}} AND {{end}}{{end}};` + "`" + `
}`)
	tmpl.Execute(&buf, t)
	return buf.String()
}
