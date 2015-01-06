package models

import (
	"database/sql"
	"time"
)

//go:generate tabler $GOFILE

// User stores user account information.
// @table
type User struct {
	ID      string    `tabler:"columnType=uuid&primary=true"`
	Email   string    `tabler:"columnType=varchar(128)"`
	Created time.Time `tabler:"columnType=timestamp"`
	db      *sql.DB   // tabler will ignore fields lacking the "tabler" tag
}

// Profile stores user attributes.
// @table
type Profile struct {
	UserID    string `tabler:"columnType=uuid&primary=true"`
	Attribute string `tabler:"columnType=varchar(64)&primary=true"`
	Value     string `tabler:"columnType=varchar(256)"`
	db        *sql.DB
}
