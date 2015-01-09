package lib

import (
	"bytes"
	"text/template"
)

// DropTableStatement ...
func (t Table) DropTableStatement() string {
	buf := bytes.Buffer{}
	tmpl := templify(`DROP TABLE {{lower .Name}};`)
	tmpl.Execute(&buf, t)
	return buf.String()
}

// DropTable ...
func (t Table) DropTable() string {
	buf := bytes.Buffer{}
	var tmpl *template.Template
	if !t.HasConn {
		tmpl = templify(`func ({{caller .Name}} {{.Name}}) DropTable() string {
    return ` + "`{{.DropTableStatement}}`" + `
}`)
	} else {
		tmpl = templify(`func ({{caller .Name}} {{.Name}}) DropTable() error {
    _, err := {{caller .Name}}.db.Exec(` + "`{{.DropTableStatement}}`)" + `
    return err
}`)
	}
	tmpl.Execute(&buf, t)
	return buf.String()
}
