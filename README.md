# tabler

Go generate syntactic sugar for SQL-backed structs

## Status 2015-01-07

This project was originally a monolithic proof-of-concept script for Go's `generate` subcommand, but has started to see some use in a side project and will remain under development.

I have not put together a formal design doc, but some enhancements on the docket include:

- Breaking the package up into a library and a main. Testing, coverage... These things are on my mind.
- Move templates out of source. *Templates may go away entirely.*
- For structs with a database connection reference (either `*sql.DB` or the very convenient [`*sqlx.DB`](https://github.com/jmoiron/sqlx)), `tabler` is going to implement ORM-like functionality. Lots of potential there, from robust queries to test fixtures.

For those looking for a `go generate` example in order to start hacking together their own utilities, take a look at [commit 698eca6ad6](https://github.com/tristanwietsma/tabler/tree/698eca6ad6d7773ccbaebcffaf45ef8629d47019).

## Introduction

Given a struct with tagged fields, `tabler` will generate methods that return strings for the following actions:
- Create Table
- Drop Table
- Insert Row
- Select Row

See the [example](https://github.com/tristanwietsma/tabler/tree/master/example) for more information.

## Installation

```bash
go get github.com/tristanwietsma/tabler
```

## Use

Add the `go:generate` directive to files with SQL-backed structs.

```go
//go:generate tabler $GOFILE
```

Add the `@table` decorator to the comment block for all target structs. Tag each field with the data type (`columnType`) and label the primary keys.

```go
// @table
type User struct {
    ID      string    `tabler:"columnType=uuid&primary=true"`
    Email   string    `tabler:"columnType=varchar(128)"`
    Created time.Time `tabler:"columnType=timestamp"`
}
```

Run `generate` and tabler will produce `*_tabler.go` files for those files containing decorated structs.

```bash
go generate
go build
```

## Tags

### Requirements

- Every field must have a `tabler` key in the tag in order to be included as a column.
- Struct fields without a `tabler` key will be ignored.
- A `columnType` attribute is required for every field.
- Every table must have at least one primary key.

### Foreign Key Convention

Fields matching the pattern `<something>ID` are assumed to be foreign keys. For example:

```go
// @table
type Profile struct {
    UserID    string `tabler:"columnType=uuid&primary=true"`
    Attribute string `tabler:"columnType=varchar(64)&primary=true"`
    Value     string `tabler:"columnType=varchar(256)"`
}
```

In the above, `UserID` will be defined as `userid uuid REFERENCES user(id)` in the table creation statement.
