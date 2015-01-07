package main

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

func parseAttr(attr string) map[string]string {
	attrMap := make(map[string]string)
	chunks := strings.Split(attr, "&")
	for _, chunk := range chunks {
		pair := strings.Split(chunk, "=")
		attrMap[pair[0]] = pair[1]
	}
	return attrMap
}

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
	c.Name = name
	attrMap := parseAttr(tag)

	// columnType
	if dtype, ok := attrMap["columnType"]; ok {
		c.Type = dtype
	} else {
		return fmt.Errorf("Missing columnType!")
	}

	// primary key
	if isP, ok := attrMap["primary"]; ok && isP == "true" {
		c.IsPrimary = true
	}

	// foreign key
	pattern := regexp.MustCompile("([A-Za-z][A-Za-z0-9]*)ID")
	match := pattern.FindStringSubmatch(name)
	if len(match) > 0 {
		c.IsForeign = true
		c.ForeignKey = fmt.Sprintf("REFERENCES %s(id)", strings.ToLower(match[1]))
	}

	return nil
}

func (c Column) String() string {
	buf := bytes.Buffer{}
	tmpl := newTmpl(`{{lower .Name}} {{.Type}}{{if .IsForeign}} {{.ForeignKey}}{{end}}`)
	tmpl.Execute(&buf, c)
	return buf.String()
}
