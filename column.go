package main

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

var (
	foreignKeyPattern = regexp.MustCompile("([A-Za-z][A-Za-z0-9]*)ID")
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
	FType      string
	Type       string
	IsPrimary  bool
	IsForeign  bool
	ForeignKey string
}

func (c *Column) init(name, ftype, tag string) error {
	c.Name = name
	c.FType = ftype
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
	match := foreignKeyPattern.FindStringSubmatch(name)
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
