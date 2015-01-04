package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"
	//"text/template"
)

// Column is a column in an SQL table.
type Column struct {
	Name       string
	Type       string
	IsPrimary  bool
	ForeignKey string
}

// Init sets the columns fields.
func (c *Column) Init(name, tag string) error {
	(*c).Name = strings.ToLower(name)

	if len((*c).Name) > 2 && (*c).Name[len((*c).Name)-2:] == "id" {
		(*c).ForeignKey = fmt.Sprintf("references %s(id)", (*c).Name[:len((*c).Name)-2])
	}

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

// Table is an SQL table.
type Table struct {
	Name    string
	Columns []Column
}

// CreateTable returns a create table statement for the table.
func (t Table) CreateTable() string {
	return "CREATE TABLE" // to do
}

// DropTable returns a drop table statement for the table.
func (t Table) DropTable() string {
	return "DROP TABLE" // to do
}

// InsertRow returns an insert row statement for the table.
func (t Table) InsertRow() string {
	return "INSERT ROW" // to do
}

// DeleteRow returns a delete row statement for the table.
func (t Table) DeleteRow() string {
	return "DELETE ROW" // to do
}

// UpdateRow returns an update row statement for the table.
func (t Table) UpdateRow() string {
	return "UPDATE ROW" // to do
}

// SelectRow returns a select row statement for the table.
func (t Table) SelectRow() string {
	return "SELECT ROW" // to do
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
				table.Name = strings.ToLower(ts.Name.Name)
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
		}
		if len(table.Columns) > 0 {
			(*i).Tables = append((*i).Tables, table)
		}
	}

	return nil
}

func (i InputFile) Write() error {
	fmt.Println(i)
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
