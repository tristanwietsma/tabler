package lib

// Table is an SQL table.
type Table struct {
	Name        string
	HasConn     bool
	Conn        string
	Columns     []Column
	PrimaryKeys []Column
}
