package main

import (
	"bytes"
	"fmt"
	"strings"
)

// Column is a column in an SQL table.
type Column struct {
	Name       string
	Type       string
	IsPrimary  bool
	IsForeign  bool
	ForeignKey string
}

// Init sets the columns fields.
func (c *Column) Init(name, tag string) error {
	(*c).Name = name

	// auto-detect foreign key
	if len(name) > 2 && name[len(name)-2:] == "ID" {
		(*c).IsForeign = true
		tbl := strings.ToLower((*c).Name[:len((*c).Name)-2])
		(*c).ForeignKey = fmt.Sprintf("REFERENCES %s(id)", tbl)
	}

	// parse attributes
	attributes := strings.Split(tag, "&")
	for _, attr := range attributes {
		pair := strings.Split(attr, "=")
		if len(pair) != 2 {
			return fmt.Errorf("Malformed tag: '%s'", attr)
		}

		switch strings.ToLower(pair[0]) {
		case "columntype":
			(*c).Type = pair[1]
		case "primary":
			if pair[1] == "true" {
				(*c).IsPrimary = true
			}
		default:
			return fmt.Errorf("Unknown attribute: '%s'", pair[0])
		}
	}

	return nil
}

func (c Column) String() string {
	buf := bytes.Buffer{}
	tmpl := newTmpl(`{{lower .Name}} {{.Type}}{{if .IsForeign}} {{.ForeignKey}}{{end}}`)
	tmpl.Execute(&buf, c)
	return buf.String()
}
