package lib

import (
	"bytes"
)

// SelectRow ...
func (t Table) SelectRow() string {
	buf := bytes.Buffer{}
	tmpl := newTmpl(
		`func ({{caller .Name}} {{.Name}}) SelectRow() string {
    return ` + "`SELECT * FROM {{lower .Name}}{{$p := len .PrimaryKeys}}{{if gt $p 0}} WHERE{{end}}{{range $j, $k := .PrimaryKeys}} {{lower $k.Name}}=?{{if lt (plus1 $j) $p}} AND{{end}}{{end}};" + "`" + `
}`)
	tmpl.Execute(&buf, t)
	return buf.String()
}
