package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var (
	filters = template.FuncMap{
		"plus1":  func(x int) int { return x + 1 },
		"lower":  func(s string) string { return strings.ToLower(s) },
		"caller": func(s string) string { return strings.ToLower(s)[0:1] },
	}
)

func newTmpl(s string) *template.Template {
	return template.Must(template.New("T").Funcs(filters).Parse(s))
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
	(*c).Name = name

	// auto-detect foreign key
	if len(name) > 2 && name[len(name)-2:] == "ID" {
		(*c).IsForeign = true
		tbl := strings.ToLower((*c).Name[:len((*c).Name)-2])
		(*c).ForeignKey = fmt.Sprintf("REFERENCES %s(id)", tbl)
	}

	// parse attributes
	attributes := strings.Split(
		strings.Trim(tag, "`"),
		",",
	)
	for _, attr := range attributes {
		pair := strings.Split(attr, ":")
		if len(pair) != 2 {
			return fmt.Errorf("Malformed tag: '%s'", attr)
		}

		switch strings.ToLower(pair[0]) {
		case "type":
			(*c).Type = pair[1]
		case "primary":
			(*c).IsPrimary = true

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
    return ` + "`" + `INSERT INTO {{lower .Name}} ({{$n := len .Columns}}{{range $i, $c := .Columns}}{{$c.Name}}{{if lt (plus1 $i) $n}}, {{end}}{{end}}) VALUES ({{range $i, $c := .Columns}}?{{if lt (plus1 $i) $n}}, {{end}}{{end}});` + "`" + `
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

// InputFile is a go file with tables.
type InputFile struct {
	PackageName string
	BuildTarget string
	Tables      []Table
}

// Init initializes an InputFile from a path.
func (i *InputFile) Init(path string) error {

	// set the build target
	if strings.HasSuffix(path, ".go") {
		root := strings.TrimSuffix(path, ".go")
		dir, file := filepath.Split(root)
		(*i).BuildTarget = filepath.Join(dir, fmt.Sprintf("%s_tabler.go", file))
	} else {
		return fmt.Errorf("File '%s' is not a Go file.", path)
	}

	f, err := parser.ParseFile(
		token.NewFileSet(),
		path,
		nil,
		parser.ParseComments,
	)
	if err != nil {
		fmt.Errorf("Unable to parse '%s': %s", path, err)
	}

	// get package name
	if f.Name != nil {
		(*i).PackageName = f.Name.Name
	} else {
		fmt.Errorf("Missing package name in '%s'", path)
	}

	// build list of tables
	var isTable bool
	for _, decl := range f.Decls {

		// get the type declaration
		tdecl, ok := decl.(*ast.GenDecl)
		if !ok || tdecl.Doc == nil {
			continue
		}

		// find the @table decorator
		isTable = false
		for _, comment := range tdecl.Doc.List {
			if strings.Contains(comment.Text, "@table") {
				isTable = true
				break
			}
		}
		if !isTable {
			continue
		}

		table := Table{}

		// get the name of the table
		for _, spec := range tdecl.Specs {
			if ts, ok := spec.(*ast.TypeSpec); ok {
				if ts.Name == nil {
					continue
				}
				table.Name = ts.Name.Name
				break
			}
		}
		if table.Name == "" {
			return fmt.Errorf("Unable to extract name from a table struct.")
		}

		// parse tags and build columns
		sdecl := tdecl.Specs[0].(*ast.TypeSpec).Type.(*ast.StructType)
		fields := sdecl.Fields.List
		for _, field := range fields {
			col := Column{}
			if err := col.Init(field.Names[0].Name, field.Tag.Value); err != nil {
				return fmt.Errorf(
					"Unable to parse tag '%s' from table '%s' in '%s': %v",
					field.Tag.Value,
					table.Name,
					path,
					err,
				)
			}
			table.Columns = append(table.Columns, col)
			if col.IsPrimary {
				table.PrimaryKeys = append(table.PrimaryKeys, col)
			}
		}

		if len(table.Columns) > 0 && len(table.PrimaryKeys) > 0 {
			(*i).Tables = append((*i).Tables, table)
		}
	}

	return nil
}

func (i InputFile) Write() error {
	buf := bytes.Buffer{}
	tmpl := newTmpl(`// generated by tabler
package {{.PackageName}}
{{range $j, $t := .Tables}}
// {{.Name}}

{{$t.CreateTable}}

{{$t.DropTable}}

{{$t.InsertRow}}

{{$t.SelectRow}}
{{end}}`)
	tmpl.Execute(&buf, i)

	fmt.Println(buf.String())
	return nil
}

func main() {
	for _, path := range os.Args[1:] {
		infile := InputFile{}
		if err := infile.Init(path); err != nil {
			log.Printf("%v", err)
			continue
		}
		infile.Write()
	}
}
