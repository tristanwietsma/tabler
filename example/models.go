package models

import (
	"database/sql"
	"time"
)

//go:generate tabler $GOFILE

// @table
type User struct {
	ID      string    `tabler:"columnType=uuid&primary=true"`
	Email   string    `tabler:"columnType=varchar(128)"`
	Created time.Time `tabler:"columnType=timestamp"`
	db      *sql.DB   // tabler will detect the database connection and build I/O methods
}

// Foreign key relationships are assumed when fields exist with names matching [a-zA-Z][a-zA-Z0-9]*ID.
// If no *sql.DB field exists, tabler will only build methods that return SQL statements.
// @table
type Profile struct {
	UserID    string `tabler:"columnType=uuid&primary=true"`
	Attribute string `tabler:"columnType=varchar(64)&primary=true"`
	Value     string `tabler:"columnType=varchar(256)"`
}

// @table
type NoPrimary struct {
	Attribute string `tabler:"columnType=varchar(64)"`
	Value     string `tabler:"columnType=varchar(256)"`
}
