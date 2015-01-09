package lib

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"
)

// Column is a column in an SQL table.
type Column struct {
	Name       string
	FType      string
	Type       string
	IsPrimary  bool
	IsForeign  bool
	ForeignKey string
}

func (c *Column) init(name, ftype, tag string) error {
	c.FType = ftype

	attribute, err := url.ParseQuery(tag)
	if err != nil {
		return fmt.Errorf("Unable to parse tag: %s", tag)
	}

	if c.Type = attribute.Get("columnType"); c.Type == "" {
		return fmt.Errorf("Missing columnType!")
	}

	if attribute.Get("primary") == "true" {
		c.IsPrimary = true
	}

	c.Name = name
	if m := foreignKeyPattern.FindStringSubmatch(name); len(m) > 0 {
		c.IsForeign = true
		c.ForeignKey = fmt.Sprintf("REFERENCES %s(id)", strings.ToLower(m[1]))
	}

	return nil
}

func (c Column) String() string {
	buf := bytes.Buffer{}
	tmpl := newTmpl(`{{lower .Name}} {{.Type}}{{if .IsForeign}} {{.ForeignKey}}{{end}}`)
	tmpl.Execute(&buf, c)
	return buf.String()
}
